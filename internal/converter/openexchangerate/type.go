package openexchangerate

type Rates map[string]float32

type ConvertResponse struct {
	Timestamp int64  `json:"timestamp"`
	Base      string `json:"base"`
	Rates     Rates  `json:"rates"`
}
