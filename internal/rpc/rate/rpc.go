package rate

import (
	"context"
	"errors"

	"github.com/satimoto/go-datastore/pkg/util"
	"github.com/satimoto/go-ferp/ferprpc"
)

func (r *RpcRateResolver) SubscribeRates(req *ferprpc.SubscribeRatesRequest, stream ferprpc.RateService_SubscribeRatesServer) error {
	cancelCtx, cancel := context.WithCancel(context.Background())
	ratesChan := r.ConverterService.SubscribeRates(cancelCtx)
	defer cancel()

	for {
		select {
		case currencyRate := <-ratesChan:
			if len(req.Currency) == 0 || req.Currency == currencyRate.Currency {
				subscribeRatesResponse := &ferprpc.SubscribeRatesResponse{
					Currency:    currencyRate.Currency,
					Rate:        currencyRate.Rate,
					RateMsat:    currencyRate.RateMsat,
					LastUpdated: currencyRate.LastUpdated.Unix(),
				}

				err := stream.Send(subscribeRatesResponse)

				if err != nil {
					util.LogOnError("FERP006", "Error sending to stream", err)
					return err
				}
			}
		case <-stream.Context().Done():
			if errors.Is(stream.Context().Err(), context.Canceled) {
				return nil
			}

			return stream.Context().Err()
		case <-r.shutdownCtx.Done():
			return nil
		}
	}
}
