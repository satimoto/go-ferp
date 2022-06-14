package rate

import (
	"github.com/satimoto/go-datastore/pkg/util"
	"github.com/satimoto/go-ferp/ferprpc"
)

func (r *RpcRateResolver) UpdateRate(response *ferprpc.SubscribeRatesResponse) {
	streamsToUpdate := r.streams

	for index, rateStream := range streamsToUpdate {
		if len(rateStream.Currency) == 0 || rateStream.Currency == response.Currency {
			err := rateStream.Stream.Send(response)

			if err != nil {
				util.LogOnError("FERP006", "Error sending to stream", err)
				r.removeRateSubscription(index)
			}
		}
	}
}
