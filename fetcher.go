package octopus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchRates(apiKey, productCode, tariffCode, rateType string) ([]Rate, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.octopus.energy/v1/products/%s/electricity-tariffs/%s/%s-unit-rates/", productCode, tariffCode, rateType)
	rates := []Rate{}

	for url != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var rateResponse RateResponse
		if err := json.Unmarshal(body, &rateResponse); err != nil {
			return nil, err
		}

		rates = append(rates, rateResponse.Results...)

		url = rateResponse.Next
	}

	return rates, nil
}
