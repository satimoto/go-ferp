package openexchangerate

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	dbUtil "github.com/satimoto/go-datastore/pkg/util"
	metrics "github.com/satimoto/go-ferp/internal/metric"
	"github.com/satimoto/go-ferp/pkg/rate"
)

const (
	openExchangeRateAPIURL = "https://openexchangerates.org"
)

type OpenExchangeRate interface {
	UpdateRates() (rate.LatestConversionRates, error)
	GetRate(currency string) (*rate.ConversionRate, error)
	GetRates() rate.LatestConversionRates
}

type OpenExchangeRateService struct {
	apiKey          string
	httpClient      *http.Client
	conversionRates rate.LatestConversionRates
	base            string
	symbols         []string
}

func NewConverter(apiKey string) OpenExchangeRate {
	return NewConverterWithClient(apiKey, http.DefaultClient)
}

func NewConverterWithClient(apiKey string, httpClient *http.Client) OpenExchangeRate {
	return &OpenExchangeRateService{
		apiKey:          apiKey,
		httpClient:      httpClient,
		conversionRates: make(rate.LatestConversionRates),
		base:            "EUR",
		symbols: []string{
			"ANG", "BAM", "BGN", "CHF", "CLP", "CZK", "DKK", "EUR", "HRK", "HUF", "ISK",
			"KRW", "LVL", "PLN", "MAD", "MKD", "NOK", "RON", "RSD", "SEK", "SGD", "THB"},
	}
}

func (s *OpenExchangeRateService) UpdateRates() (rate.LatestConversionRates, error) {
	latestConversionRates := make(rate.LatestConversionRates)

	conversionRates, err := s.queryRates(s.symbols)

	if err != nil {
		return nil, err
	}

	for _, conversionRate := range conversionRates {
		latestConversionRates[conversionRate.ToCurrency] = conversionRate
		s.conversionRates[conversionRate.ToCurrency] = conversionRate
	}

	return latestConversionRates, nil
}

func (s *OpenExchangeRateService) queryRates(symbols []string) ([]rate.ConversionRate, error) {
	conversionRates := []rate.ConversionRate{}

	if len(symbols) > 0 {
		values := url.Values{}
		values.Set("app_id", s.apiKey)
		values.Set("symbols", strings.Join(symbols, ","))

		requestUrl := fmt.Sprintf("%s/api/latest.json?%s", openExchangeRateAPIURL, values.Encode())
		request, err := http.NewRequest(http.MethodGet, requestUrl, nil)

		if err != nil {
			metrics.RecordError("FERP014", "Error forming request", err)
			log.Printf("FERP014: Url=%v", requestUrl)
			return nil, errors.New("error forming request")
		}

		response, err := s.httpClient.Do(request)

		if err != nil {
			metrics.RecordError("FERP015", "Error making request", err)
			dbUtil.LogHttpRequest("FERP015", requestUrl, request, false)
			return nil, errors.New("error making request")
		}

		convertResponse, err := UnmarshalConvertResponse(response.Body)

		if err != nil {
			metrics.RecordError("FERP016", "Error unmarshalling response", err)
			dbUtil.LogHttpResponse("FERP016", requestUrl, response, false)
			return nil, errors.New("error unmarshalling response")
		}

		var baseRate float32 = 1.0
		var ok bool = true

		if s.base != convertResponse.Base {
			baseRate, ok = convertResponse.Rates[s.base]
		}

		if ok {
			for symbol, value := range convertResponse.Rates {
				if symbol != s.base {
					convertedValue := value / baseRate
					conversionRates = append(conversionRates, rate.ConversionRate{
						FromCurrency: "EUR",
						ToCurrency:   symbol,
						Rate:         convertedValue,
						LastUpdated:  time.Now(),
					})
				}
			}
		}
	}

	return conversionRates, nil
}

func (s *OpenExchangeRateService) GetRate(currency string) (*rate.ConversionRate, error) {
	if conversionRate, ok := s.conversionRates[currency]; ok {
		return &conversionRate, nil
	}

	return nil, errors.New("no conversion rate available")
}

func (s *OpenExchangeRateService) GetRates() rate.LatestConversionRates {
	return s.conversionRates
}
