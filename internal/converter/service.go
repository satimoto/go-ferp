package converter

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/satimoto/go-ferp/internal/converter/currencyconverter"
	"github.com/satimoto/go-ferp/internal/exchange"
	"github.com/satimoto/go-ferp/internal/rpc"
	"github.com/satimoto/go-ferp/pkg/rate"
)

type Converter interface {
	Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup)
}

type ConverterService struct {
	ExchangeService         exchange.Exchange
	RpcService              rpc.Rpc
	currencyConverterClient currencyconverter.CurrencyConverter
	conversionRates         rate.LatestConversionRates
}

func NewService(exchangeService exchange.Exchange, rpcService rpc.Rpc) Converter {
	return &ConverterService{
		ExchangeService:         exchangeService,
		RpcService:              rpcService,
		currencyConverterClient: currencyconverter.NewConverter(os.Getenv("CURRENCY_CONVERTER_API_KEY")),
		conversionRates:         make(rate.LatestConversionRates),
	}
}

func (s *ConverterService) Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	s.ExchangeService.AddRateHandler(s.handleRate)

	go s.startUpdateLoop(shutdownCtx, waitGroup)
}

func (s *ConverterService) handleRate(currency string, currencyRate rate.CurrencyRate) {
	rateService := s.RpcService.GetRateService()

	for _, conversionRate := range s.conversionRates {
		if currency == conversionRate.FromCurrency {
			convertedCurrencyRate := rate.CurrencyRate{
				Rate:        int64(float32(currencyRate.Rate) * conversionRate.Rate),
				RateMsat:    int64(float32(currencyRate.RateMsat) * conversionRate.Rate),
				LastUpdated: currencyRate.LastUpdated,
			}

			log.Printf("%s: %v sats / %v millisats", conversionRate.ToCurrency, convertedCurrencyRate.Rate, convertedCurrencyRate.RateMsat)
			rateService.UpdateRate(NewSubscribeRatesResponse(conversionRate.ToCurrency, convertedCurrencyRate, conversionRate))
		}
	}
}

func (s *ConverterService) startUpdateLoop(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	log.Printf("Starting Converter service")
	waitGroup.Add(1)

	time.Now().Minute()

updateLoop:
	for {
		currencyRates, err := s.currencyConverterClient.UpdateRates()

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
