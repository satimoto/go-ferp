package kraken

type Ticker struct {
	Ask                        []string `json:"a"`
	Bid                        []string `json:"b"`
	Last                       []string `json:"c"`
	Volume                     []string `json:"v"`
	VolumeWeightedAveragePrice []string `json:"p"`
	Trades                     []int64  `json:"t"`
	Low                        []string `json:"l"`
	High                       []string `json:"h"`
	Open                       string   `json:"o"`
}

type TickerResponse struct {
	Error []interface{}     `json:"error"`
	Data  map[string]Ticker `json:"result"`
}
