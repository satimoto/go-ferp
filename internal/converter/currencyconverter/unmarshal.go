package currencyconverter

import (
	"encoding/json"
	"io"
)

func UnmarshalConvertResponse(body io.ReadCloser) (ConvertResponse, error) {
	var response ConvertResponse

	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
