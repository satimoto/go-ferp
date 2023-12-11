package converter

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/satimoto/go-ferp/internal/converter/currencyconverter"
	"github.com/satimoto/go-ferp/internal/converter/openexchangerate"
	"github.com/satimoto/go-ferp/internal/exchange"
	"github.com/satimoto/go-ferp/pkg/rate"
)

type Converter interface {
	Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup)
	SubscribeRates(cancelCtx context.Context) chan *rate.CurrencyRate
}

type ConverterService struct {
	ExchangeService         exchange.Exchange
	currencyConverterClient currencyconverter.CurrencyConverter
	openExchangeRateClient  openexchangerate.OpenExchangeRate
	conversionRates         rate.LatestConversionRates
	currencyRates           rate.LatestCurrencyRates
	rateSubscriptions       map[string]chan *rate.CurrencyRate
}

func NewService(exchangeService exchange.Exchange) Converter {
	return &ConverterService{
		ExchangeService:         exchangeService,
		currencyConverterClient: currencyconverter.NewConverter(os.Getenv("CURRENCY_CONVERTER_API_KEY")),
		openExchangeRateClient:  openexchangerate.NewConverter(os.Getenv("OPEN_EXCHANGE_RATE_API_KEY")),
		conversionRates:         make(rate.LatestConversionRates),
		currencyRates:           make(rate.LatestCurrencyRates),
		rateSubscriptions:       make(map[string]chan *rate.CurrencyRate),
	}
}

func (s *ConverterService) Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	go s.startSubscriptionListener(shutdownCtx)
	go s.startUpdateLoop(shutdownCtx, waitGroup)
}

func (s *ConverterService) SubscribeRates(cancelCtx context.Context) chan *rate.CurrencyRate {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	s.rateSubscriptions[id] = make(chan *rate.CurrencyRate)

	go s.sendInitialRates(s.rateSubscriptions[id])
	go s.waitForSubscriptionCancellation(cancelCtx, id)

	return s.rateSubscriptions[id]
}

func (s *ConverterService) handleCurrencyRate(currencyRate *rate.CurrencyRate) {
	convertedCurrencyRate := &rate.CurrencyRate{
		Currency:    currencyRate.Currency,
		Rate:        currencyRate.Rate,
		RateMsat:    currencyRate.RateMsat,
		LastUpdated: currencyRate.LastUpdated,
	}

	metricCurrencyRateSatoshis.WithLabelValues(convertedCurrencyRate.Currency).Set(float64(currencyRate.Rate))

	s.currencyRates[currencyRate.Currency] = *convertedCurrencyRate
	s.updateRateSubscriptions(convertedCurrencyRate)

	// Update conversion rates
	for _, conversionRate := range s.conversionRates {
		if currencyRate.Currency == conversionRate.FromCurrency {
			rateMsat := int64(float32(currencyRate.RateMsat) / conversionRate.Rate)
			rateSat := rateMsat / 1000
			convertedCurrencyRate := &rate.CurrencyRate{
				Currency:    conversionRate.ToCurrency,
				Rate:        rateSat,
				RateMsat:    rateMsat,
				LastUpdated: currencyRate.LastUpdated,
			}

			metricCurrencyRateSatoshis.WithLabelValues(conversionRate.ToCurrency).Set(float64(rateSat))

			s.currencyRates[conversionRate.ToCurrency] = *convertedCurrencyRate
			s.updateRateSubscriptions(convertedCurrencyRate)
		}
	}
}

func (s *ConverterService) sendInitialRates(ratesChan chan *rate.CurrencyRate) {
	<-time.After(100 * time.Millisecond)

	for _, currencyRate := range s.currencyRates {
		ratesChan <- &currencyRate
	}
}

func (s *ConverterService) startSubscriptionListener(shutdownCtx context.Context) {
	ratesChan := s.ExchangeService.SubscribeRates(shutdownCtx)

updateLoop:
	for {
		select {
		case currencyRate := <-ratesChan:
			s.updateRateSubscriptions(currencyRate)
			s.handleCurrencyRate(currencyRate)
		case <-shutdownCtx.Done():
			break updateLoop
		}
	}

}

func (s *ConverterService) startUpdateLoop(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	log.Printf("Starting Converter service")
	waitGroup.Add(1)

updateLoop:
	for {
		currencyRates, err := s.currencyConverterClient.UpdateRates()

		if err != nil {
			log.Printf("Trying OpenExchangeRate client")
			currencyRates, err = s.openExchangeRateClient.UpdateRates()
		}

		if err != nil {
			log.Printf("Using backup rates")
			currencyRates = s.backupRates()
		}

		for currency, currencyRate := range currencyRates {
			s.conversionRates[currency] = currencyRate
		}

		select {
		case <-shutdownCtx.Done():
			log.Printf("Shutting down Converter service")
			break updateLoop
		case <-time.After(time.Duration(61-time.Now().Minute()) * time.Minute):
		}
	}

	log.Printf("Converter service shut down")
	waitGroup.Done()
}

func (s *ConverterService) backupRates() rate.LatestConversionRates {
	updatedAt := time.Now()
	backupConversionRates := make(rate.LatestConversionRates)
	backupConversionRates["ANG"] = s.backupRate("EUR", "ANG", 1.93, updatedAt)
	backupConversionRates["BAM"] = s.backupRate("EUR", "BAM", 1.96, updatedAt)
	backupConversionRates["BGN"] = s.backupRate("EUR", "BGN", 1.96, updatedAt)
	backupConversionRates["CHF"] = s.backupRate("EUR", "CHF", 0.95, updatedAt)
	backupConversionRates["CLP"] = s.backupRate("EUR", "CLP", 946.06, updatedAt)
	backupConversionRates["CZK"] = s.backupRate("EUR", "CZK", 24.35, updatedAt)
	backupConversionRates["DKK"] = s.backupRate("EUR", "DKK", 7.46, updatedAt)
	backupConversionRates["HRK"] = s.backupRate("EUR", "HRK", 7.53, updatedAt)
	backupConversionRates["HUF"] = s.backupRate("EUR", "HUF", 381.61, updatedAt)
	backupConversionRates["ISK"] = s.backupRate("EUR", "ISK", 150.71, updatedAt)
	backupConversionRates["KRW"] = s.backupRate("EUR", "KRW", 1419.87, updatedAt)
	backupConversionRates["LVL"] = s.backupRate("EUR", "LVL", 0.70, updatedAt)
	backupConversionRates["PLN"] = s.backupRate("EUR", "PLN", 4.33, updatedAt)
	backupConversionRates["MAD"] = s.backupRate("EUR", "MAD", 10.94, updatedAt)
	backupConversionRates["MKD"] = s.backupRate("EUR", "MKD", 61.53, updatedAt)
	backupConversionRates["NOK"] = s.backupRate("EUR", "NOK", 11.79, updatedAt)
	backupConversionRates["RON"] = s.backupRate("EUR", "RON", 4.97, updatedAt)
	backupConversionRates["RSD"] = s.backupRate("EUR", "RSD", 117.17, updatedAt)
	backupConversionRates["SEK"] = s.backupRate("EUR", "SEK", 11.28, updatedAt)
	backupConversionRates["SGD"] = s.backupRate("EUR", "SGD", 1.45, updatedAt)
	backupConversionRates["THB"] = s.backupRate("EUR", "THB", 38.40, updatedAt)
	backupConversionRates["GBP"] = s.backupRate("EUR", "GBP", 0.86, updatedAt)
	backupConversionRates["USD"] = s.backupRate("EUR", "USD", 10.08, updatedAt)

	return backupConversionRates
}

func (s *ConverterService) backupRate(from, to string, value float32, updatedAt time.Time) rate.ConversionRate {
	return rate.ConversionRate{
		FromCurrency: from,
		ToCurrency: to,
		Rate: value,
		LastUpdated: updatedAt,
	}
}



func (s *ConverterService) updateRateSubscriptions(currencyRate *rate.CurrencyRate) {
	for _, rateSubscription := range s.rateSubscriptions {
		rateSubscription <- currencyRate
	}
}

func (s *ConverterService) waitForSubscriptionCancellation(cancelCtx context.Context, id string) {
	<-cancelCtx.Done()
	close(s.rateSubscriptions[id])
	delete(s.rateSubscriptions, id)
}
