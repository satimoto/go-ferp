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
	rateSubscriptions       map[string]chan *rate.CurrencyRate
}

func NewService(exchangeService exchange.Exchange) Converter {
	return &ConverterService{
		ExchangeService:         exchangeService,
		currencyConverterClient: currencyconverter.NewConverter(os.Getenv("CURRENCY_CONVERTER_API_KEY")),
		openExchangeRateClient:  openexchangerate.NewConverter(os.Getenv("OPEN_EXCHANGE_RATE_API_KEY")),
		conversionRates:         make(rate.LatestConversionRates),
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

			s.updateRateSubscriptions(convertedCurrencyRate)
		}
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
			log.Printf("Using OpenExchangeRate client")
			currencyRates, err = s.openExchangeRateClient.UpdateRates()
		}

		if err == nil {
			for currency, currencyRate := range currencyRates {
				s.conversionRates[currency] = currencyRate
			}
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
