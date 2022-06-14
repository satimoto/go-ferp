package rate

import "time"

type ConversionRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float32
	LastUpdated  time.Time
}

type LatestConversionRates map[string]ConversionRate

type CurrencyRate struct {
	Rate        int64
	RateMsat    int64
	LastUpdated time.Time
}

type LatestCurrencyRates map[string]CurrencyRate
