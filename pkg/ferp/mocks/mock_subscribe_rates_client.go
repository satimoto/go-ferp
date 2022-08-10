package mocks

import (
	"github.com/satimoto/go-ferp/ferprpc"
	"google.golang.org/grpc"
)

type MockSubscribeRatesClient struct {
	grpc.ClientStream
	recvChan <-chan *ferprpc.SubscribeRatesResponse
}

func NewMockSubscribeRatesClient(recvChan <-chan *ferprpc.SubscribeRatesResponse) ferprpc.RateService_SubscribeRatesClient {
	clientStream := NewMockClientStream()
	return &MockSubscribeRatesClient{
		ClientStream: clientStream,
		recvChan: recvChan,
	}
}

func (c *MockSubscribeRatesClient) Recv() (*ferprpc.SubscribeRatesResponse, error) {
	receive := <-c.recvChan
	return receive, nil
}

