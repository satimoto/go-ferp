package kraken

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/satimoto/go-datastore/pkg/util"
	metrics "github.com/satimoto/go-ferp/internal/metric"
	"github.com/satimoto/go-ferp/pkg/rate"
)

const (
	krakenAPIURL  = "https://api.kraken.com"
	krakenTicker  = "Ticker"
	krakenVersion = "0"
)

type Kraken interface {
	UpdateRates() (rate.LatestCurrencyRates, error)
	GetRate(currency string) (*rate.CurrencyRate, error)
	GetRates() rate.LatestCurrencyRates
}

type KrakenExchange struct {
	httpClient    *http.Client
	currencyRates rate.LatestCurrencyRates
}

func NewExchange() Kraken {
	return &KrakenExchange{
		httpClient:    http.DefaultClient,
		currencyRates: make(rate.LatestCurrencyRates),
	}
}

func NewExchangeWithClient(httpClient *http.Client) Kraken {
	return &KrakenExchange{
		httpClient:    httpClient,
		currencyRates: make(rate.LatestCurrencyRates),
	}
}

func (e *KrakenExchange) UpdateRates() (rate.LatestCurrencyRates, error) {
	values := url.Values{}
	values.Set("pair", "XBTEUR,XBTGBP,XBTUSD")
	values.Set("interval", "1")

	requestUrl := fmt.Sprintf("%s/%s/public/%s?%s", krakenAPIURL, krakenVersion, krakenTicker, values.Encode())
	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)

	if err != nil {
		metrics.RecordError("FERP007", "Error forming request", err)
		log.Printf("FERP007: Url=%v", requestUrl)
		return nil, errors.New("error forming request")
	}

	response, err := e.httpClient.Do(request)

	if err != nil {
		metrics.RecordError("FERP008", "Error making request", err)
		util.LogHttpRequest("FERP008", requestUrl, request, false)
		return nil, errors.New("error making request")
	}

	tickerResponse, err := UnmarshalTickerResponse(response.Body)

	if err != nil {
		metrics.RecordError("FERP009", "Error unmarshalling response", err)
		util.LogHttpResponse("FERP009", requestUrl, response, false)
		return nil, errors.New("error unmarshalling response")
	}

	for pair, value := range tickerResponse.Data {
		currency := pair[len(pair)-3:]
		price, err := strconv.ParseFloat(value.Last[0], 64)

		if err != nil {
			metrics.RecordError("FERP010", "Error parsing float", err)
			log.Printf("FERP010: Value=%v", value.Last[0])
			continue
		}

		currencyRate := rate.CurrencyRate{
			Currency:    currency,
			Rate:        int64(100_000_000 / price),
			RateMsat:    int64(100_000_000_000 / price),
			LastUpdated: time.Now(),
		}

		e.currencyRates[currency] = currencyRate
	}

	return e.currencyRates, nil
}

func (e *KrakenExchange) GetRate(currency string) (*rate.CurrencyRate, error) {
	if currencyRate, ok := e.currencyRates[currency]; ok {
		return &currencyRate, nil
	}

	return nil, errors.New("no currency rate available")
}

func (e *KrakenExchange) GetRates() rate.LatestCurrencyRates {
	return e.currencyRates
}
