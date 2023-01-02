package converter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricCurrencyRateSatoshis = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ferp_currency_rate_satoshis",
		Help: "The currency exchange rate in satoshis",
	}, []string{"currency"})
)
