package currencyconverter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	dbUtil "github.com/satimoto/go-datastore/pkg/util"
	"github.com/satimoto/go-ferp/internal/util"
	"github.com/satimoto/go-ferp/pkg/rate"
)

const (
	currencyConverterAPIURL  = "https://free.currconv.com"
	currencyConverterVersion = "v7"
)

type CurrencyConverter interface {
	UpdateRates() (rate.LatestConversionRates, error)
	GetRate(currency string) (*rate.ConversionRate, error)
	GetRates() rate.LatestConversionRates
}

type CurrencyConverterService struct {
	apiKey          string
	httpClient      *http.Client
	conversionRates rate.LatestConversionRates
	pairs           []string
}

func NewConverter(apiKey string) CurrencyConverter {
	return NewConverterWithClient(apiKey, http.DefaultClient)
}

func NewConverterWithClient(apiKey string, httpClient *http.Client) CurrencyConverter {
	return &CurrencyConverterService{
		apiKey:          apiKey,
		httpClient:      httpClient,
		conversionRates: make(rate.LatestConversionRates),
		pairs: []string{
			"EUR_ANG", "EUR_BAM", "EUR_BGN", "EUR_CHF", "EUR_CLP",
			"EUR_CZK", "EUR_DKK", "EUR_HRK", "EUR_HUF", "EUR_ISK",
			"EUR_KRW", "EUR_LVL", "EUR_PLN", "EUR_MAD", "EUR_MKD",
			"EUR_NOK", "EUR_RON", "EUR_RSD", "EUR_SEK", "EUR_SGD",
			"EUR_THB"},
	}
}

func (s *CurrencyConverterService) UpdateRates() (rate.LatestConversionRates, error) {
	latestConversionRates := make(rate.LatestConversionRates)
	pairCount := len(s.pairs) / 2

	for i := 0; i <= pairCount; i++ {
		startIndex := i * 2
		endIndex := util.MinInt(len(s.pairs), startIndex+2)
		conversionRates, err := s.queryRates(s.pairs[startIndex:endIndex])

		if err == nil {
			for _, conversionRate := range conversionRates {
				latestConversionRates[conversionRate.ToCurrency] = conversionRate
				s.conversionRates[conversionRate.ToCurrency] = conversionRate
			}
		}
	}

	return latestConversionRates, nil
}

func (s *CurrencyConverterService) queryRates(pairs []string) ([]rate.ConversionRate, error) {
	conversionRates := []rate.ConversionRate{}

	if len(pairs) > 0 {
		values := url.Values{}
		values.Set("q", strings.Join(pairs, ","))
		values.Set("compact", "ultra")
		values.Set("apiKey", s.apiKey)

		requestUrl := fmt.Sprintf("%s/api/%s/convert?%s", currencyConverterAPIURL, currencyConverterVersion, values.Encode())
		request, err := http.NewRequest(http.MethodGet, requestUrl, nil)

		if err != nil {
			dbUtil.LogOnError("FERP011", "Error forming request", err)
			log.Printf("FERP011: Url=%v", requestUrl)
			return nil, errors.New("error forming request")
		}

		response, err := s.httpClient.Do(request)

		if err != nil {
			dbUtil.LogOnError("FERP012", "Error making request", err)
			dbUtil.LogHttpRequest("FERP012", requestUrl, request, false)
			return nil, errors.New("error making request")
		}

		convertResponse, err := UnmarshalConvertResponse(response.Body)

		if err != nil {
			dbUtil.LogOnError("FERP013", "Error unmarshalling response", err)
			dbUtil.LogHttpResponse("FERP013", requestUrl, response, false)
			return nil, errors.New("error unmarshalling response")
		}

		for pair, value := range convertResponse {
			fromCurrency := pair[:3]
			toCurrency := pair[len(pair)-3:]

			conversionRates = append(conversionRates, rate.ConversionRate{
				FromCurrency: fromCurrency,
				ToCurrency:   toCurrency,
				Rate:         value,
				LastUpdated:  time.Now(),
			})

			log.Printf("1 %s = %v %s", fromCurrency, value, toCurrency)
		}
	}

	return conversionRates, nil
}

func (s *CurrencyConverterService) GetRate(currency string) (*rate.ConversionRate, error) {
	if conversionRate, ok := s.conversionRates[currency]; ok {
		return &conversionRate, nil
	}

	return nil, errors.New("no conversion rate available")
}

func (s *CurrencyConverterService) GetRates() rate.LatestConversionRates {
	return s.conversionRates
}
