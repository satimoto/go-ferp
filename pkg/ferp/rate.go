package ferp

import (
	"context"

	"github.com/satimoto/go-ferp/ferprpc"
	"google.golang.org/grpc"
)

func (s *FerpService) SubscribeRates(ctx context.Context, in *ferprpc.SubscribeRatesRequest, opts ...grpc.CallOption) (ferprpc.RateService_SubscribeRatesClient, error) {
	return s.getRateClient().SubscribeRates(ctx, in, opts...)
}

func (s *FerpService) getRateClient() ferprpc.RateServiceClient {
	if s.rateClient == nil {
		client := ferprpc.NewRateServiceClient(s.clientConn)
		s.rateClient = &client
	}

	return *s.rateClient
}
