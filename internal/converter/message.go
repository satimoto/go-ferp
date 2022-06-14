package converter

import (
	"github.com/satimoto/go-ferp/ferprpc"
	"github.com/satimoto/go-ferp/pkg/rate"
)

func NewSubscribeRatesResponse(currency string, currencyRate rate.CurrencyRate, conversionRate rate.ConversionRate) *ferprpc.SubscribeRatesResponse {
	return &ferprpc.SubscribeRatesResponse{
		Currency: currency,
		Rate:     currencyRate.Rate,
		RateMsat: currencyRate.RateMsat,
		ConversionRate: &ferprpc.ConversionRate{
			Currency:    conversionRate.FromCurrency,
			Rate:        conversionRate.Rate,
			LastUpdated: conversionRate.LastUpdated.Unix(),
		},
		LastUpdated: currencyRate.LastUpdated.Unix(),
	}
}
