package exchange

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/satimoto/go-ferp/internal/exchange/kraken"
	"github.com/satimoto/go-ferp/pkg/rate"
)

type RateHandler func(currency string, currencyRate rate.CurrencyRate)

type Exchange interface {
	GetRate(currency string) (*rate.CurrencyRate, error)
	Start(shutdownCtx context.Context, waitGroup *sync.WaitGroup)
	SubscribeRates(cancelCtx context.Context) chan *rate.CurrencyRate
}

type ExchangeService struct {
	krakenClient      kraken.Kraken
	rateSubscriptions map[string]chan *rate.CurrencyRate
}

func NewService() Exchange {
	return &ExchangeService{
		krakenClient:      kraken.NewExchange(),
		rateSubscriptions: make(map[string]chan *rate.CurrencyRate),
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

func (s *ExchangeService) SubscribeRates(cancelCtx context.Context) chan *rate.CurrencyRate {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	s.rateSubscriptions[id] = make(chan *rate.CurrencyRate)

	go s.waitForSubscriptionCancellation(cancelCtx, id)

	return s.rateSubscriptions[id]
}

func (s *ExchangeService) startUpdateLoop(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	log.Printf("Starting Exchange service")
	waitGroup.Add(1)

updateLoop:
	for {
		currencyRates, err := s.krakenClient.UpdateRates()

		if err == nil {
			for _, currencyRate := range currencyRates {
				s.updateRateSubscriptions(currencyRate)
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

func (s *ExchangeService) updateRateSubscriptions(currencyRate rate.CurrencyRate) {
	for _, rateSubscription := range s.rateSubscriptions {
		rateSubscription <- &currencyRate
	}
}

func (s *ExchangeService) waitForSubscriptionCancellation(cancelCtx context.Context, id string) {
	<-cancelCtx.Done()
	close(s.rateSubscriptions[id])
	delete(s.rateSubscriptions, id)
}
