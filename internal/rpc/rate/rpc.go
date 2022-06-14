package rate

import (
	"github.com/satimoto/go-ferp/ferprpc"
)

func (r *RpcRateResolver) SubscribeRates(req *ferprpc.SubscribeRatesRequest, stream ferprpc.RateService_SubscribeRatesServer) error {
	r.streams = append(r.streams, &RateStream{
		Currency: req.Currency,
		Stream:   stream,
	})

	return nil
}

func (r *RpcRateResolver) removeRateSubscription(index int) {
	r.streams = append(r.streams[:index], r.streams[index+1:]...)
}
