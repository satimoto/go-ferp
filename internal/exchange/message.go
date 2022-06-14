package exchange

import (
	"github.com/satimoto/go-ferp/ferprpc"
	"github.com/satimoto/go-ferp/pkg/rate"
)

func NewSubscribeRatesResponse(currency string, currencyRate rate.CurrencyRate) *ferprpc.SubscribeRatesResponse {
	return &ferprpc.SubscribeRatesResponse{
		Currency:    currency,
		Rate:        currencyRate.Rate,
		RateMsat:    currencyRate.RateMsat,
		LastUpdated: currencyRate.LastUpdated.Unix(),
	}
}
