package mocks

import (
	"errors"

	"github.com/satimoto/go-ferp/ferprpc"
	"google.golang.org/grpc"
)

func (s *MockFerpService) SubscribeRates(in *ferprpc.SubscribeRatesRequest, opts ...grpc.CallOption) (ferprpc.RateService_SubscribeRatesClient, error) {
	if len(s.subscribeRatesMockData) == 0 {
		return nil, errors.New("NotFound")
	}

	response := s.subscribeRatesMockData[0]
	s.subscribeRatesMockData = s.subscribeRatesMockData[1:]
	return response, nil
}

func (s *MockFerpService) NewSubscribeTransactionsMockData() (chan<- *ferprpc.SubscribeRatesResponse) {
	recvChan := make(chan *ferprpc.SubscribeRatesResponse)
	s.subscribeRatesMockData = append(s.subscribeRatesMockData, NewMockSubscribeRatesClient(recvChan))

	return recvChan
}