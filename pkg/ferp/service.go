package ferp

import (
	"context"

	"github.com/satimoto/go-datastore/pkg/util"
	"github.com/satimoto/go-ferp/ferprpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Ferp interface {
	SubscribeRates(ctx context.Context, in *ferprpc.SubscribeRatesRequest, opts ...grpc.CallOption) (ferprpc.RateService_SubscribeRatesClient, error)
}

type FerpService struct {
	clientConn       *grpc.ClientConn
	rateClient       *ferprpc.RateServiceClient
}

func NewService(address string) Ferp {
	clientConn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	util.PanicOnError("FERP005", "Error connecting to FERP RPC address", err)

	return &FerpService{
		clientConn: clientConn,
	}
}
