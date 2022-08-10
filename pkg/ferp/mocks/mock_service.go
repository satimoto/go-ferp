package mocks

import (
	"github.com/satimoto/go-ferp/ferprpc"
)

type MockFerpService struct {
	subscribeRatesMockData []ferprpc.RateService_SubscribeRatesClient
}

func NewService() *MockFerpService {
	return &MockFerpService{}
}
