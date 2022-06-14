package rate

import "github.com/satimoto/go-ferp/ferprpc"

type RateStream struct {
	Currency string
	Stream   ferprpc.RateService_SubscribeRatesServer
}

type RpcRateResolver struct {
	streams []*RateStream
}

func NewResolver() *RpcRateResolver {
	return &RpcRateResolver{}
}
