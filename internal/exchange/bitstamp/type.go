package bitstamp

type TickerResponse struct {
	Ask                        string `json:"ask"`
	Bid                        string `json:"bid"`
	Last                       string `json:"last"`
	Volume                     string `json:"volume"`
	VolumeWeightedAveragePrice string `json:"vwap"`
	Low                        string `json:"low"`
	High                       string `json:"high"`
	Open                       string `json:"open"`
}
