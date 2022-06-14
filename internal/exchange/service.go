package exchange

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/satimoto/go-ferp/internal/exchange/kraken"
	"github.com/satimoto/go-ferp/internal/rate"
	"github.com/satimoto/go-ferp/internal/rpc"
)

type RateHandler func(currency string, currencyRate rate.CurrencyRate)

type Exchange interface {
	AddRateHandler(handler RateHandler)
	GetRate(currency string) (*rate.CurrencyRate, error)
	Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup)
}

type ExchangeService struct {
	RpcService   rpc.Rpc
	krakenClient kraken.Kraken
	rateHandlers []RateHandler
}

func NewService(rpcService rpc.Rpc) Exchange {
	return &ExchangeService{
		RpcService:   rpcService,
		krakenClient: kraken.NewExchange(),
	}
}

func (s *ExchangeService) AddRateHandler(handler RateHandler) {
	s.rateHandlers = append(s.rateHandlers, handler)

	for currency, currencyRate := range s.krakenClient.GetRates() {
		handler(currency, currencyRate)
	}
}

func (s *ExchangeService) GetRate(currency string) (*rate.CurrencyRate, error) {
	currencyRate, err := s.krakenClient.GetRate(currency)

	if err != nil {
		return nil, err
	}

	return currencyRate, nil
}

func (s *ExchangeService) Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	go s.startUpdateLoop(shutdownCtx, waitGroup)
}

func (s *ExchangeService) startUpdateLoop(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	log.Printf("Starting Exchange service")
	waitGroup.Add(1)

	rateService := s.RpcService.GetRateService()

updateLoop:
	for {
		currencyRates, err := s.krakenClient.UpdateRates()

		if err == nil {
			for currency, currencyRate := range currencyRates {
				rateService.UpdateRate(NewSubscribeRatesResponse(currency, currencyRate))

				s.updateRateHandlers(currency, currencyRate)
			}
		}

		select {
		case <-shutdownCtx.Done():
			log.Printf("Shutting down Exchange service")
			break updateLoop
		case <-time.After(1 * time.Minute):
		}
	}

	log.Printf("Exchange service shut down")
	waitGroup.Done()
}

func (s *ExchangeService) updateRateHandlers(currency string, currencyRate rate.CurrencyRate) {
	for _, rateHandler := range s.rateHandlers {
		rateHandler(currency, currencyRate)
	}
}
