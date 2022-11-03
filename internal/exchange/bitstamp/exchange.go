package bitstamp

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/satimoto/go-datastore/pkg/util"
	metrics "github.com/satimoto/go-ferp/internal/metric"
	"github.com/satimoto/go-ferp/pkg/rate"
)

const (
	bitstampAPIURL  = "https://www.bitstamp.net/api"
	bitstampTicker  = "ticker"
	bitstampVersion = "v2"
)

type Bitstamp interface {
	UpdateRates() (rate.LatestCurrencyRates, error)
	GetRate(currency string) (*rate.CurrencyRate, error)
	GetRates() rate.LatestCurrencyRates
}

type BitstampExchange struct {
	httpClient    *http.Client
	currencyRates rate.LatestCurrencyRates
	pairs         []string
}

func NewExchange() Bitstamp {
	return &BitstampExchange{
		httpClient:    http.DefaultClient,
		currencyRates: make(rate.LatestCurrencyRates),
		pairs:         []string{"btceur", "btcgbp", "btcusd"},
	}
}

func NewExchangeWithClient(httpClient *http.Client) Bitstamp {
	return &BitstampExchange{
		httpClient:    httpClient,
		currencyRates: make(rate.LatestCurrencyRates),
	}
}

func (e *BitstampExchange) UpdateRates() (rate.LatestCurrencyRates, error) {
	for _, pair := range e.pairs {
		currencyRate, err := e.queryRates(pair)

		if err == nil {
			e.currencyRates[currencyRate.Currency] = *currencyRate
		}
	}

	return e.currencyRates, nil
}

func (e *BitstampExchange) queryRates(pair string) (*rate.CurrencyRate, error) {
	requestUrl := fmt.Sprintf("%s/%s/%s/%s", bitstampAPIURL, bitstampVersion, bitstampTicker, pair)
	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)

	if err != nil {
		metrics.RecordError("FERP017", "Error forming request", err)
		log.Printf("FERP017: Url=%v", requestUrl)
		return nil, errors.New("error forming request")
	}

	response, err := e.httpClient.Do(request)

	if err != nil {
		metrics.RecordError("FERP018", "Error making request", err)
		util.LogHttpRequest("FERP018", requestUrl, request, false)
		return nil, errors.New("error making request")
	}

	tickerResponse, err := UnmarshalTickerResponse(response.Body)

	if err != nil {
		metrics.RecordError("FERP019", "Error unmarshalling response", err)
		util.LogHttpResponse("FERP019", requestUrl, response, false)
		return nil, errors.New("error unmarshalling response")
	}

	price, err := strconv.ParseFloat(tickerResponse.Last, 64)

	if err != nil {
		metrics.RecordError("FERP020", "Error parsing float", err)
		log.Printf("FERP020: Value=%v", tickerResponse.Last)
		return nil, errors.New("error parsing float")
	}

	currency := strings.ToUpper(pair[len(pair)-3:])
	currencyRate := rate.CurrencyRate{
		Currency:    currency,
		Rate:        int64(100_000_000 / price),
		RateMsat:    int64(100_000_000_000 / price),
		LastUpdated: time.Now(),
	}

	return &currencyRate, nil
}

func (e *BitstampExchange) GetRate(currency string) (*rate.CurrencyRate, error) {
	if currencyRate, ok := e.currencyRates[currency]; ok {
		return &currencyRate, nil
	}

	return nil, errors.New("no currency rate available")
}

func (e *BitstampExchange) GetRates() rate.LatestCurrencyRates {
	return e.currencyRates
}
