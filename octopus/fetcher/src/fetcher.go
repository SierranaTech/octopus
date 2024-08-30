package octopus

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Rate struct {
	ValueIncVAT float64   `json:"value_inc_vat"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
}

type RateResponse struct {
	Results []Rate `json:"results"`
	Next    string `json:"next"`
}

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

		body, err := ioutil.ReadAll(resp.Body)
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

func WriteRatesToInflux(rates []Rate, rateType string, influxUrl, token, org, bucket string) error {
	client := influxdb2.NewClient(influxUrl, token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(org, bucket)

	for _, rate := range rates {
		p := influxdb2.NewPointWithMeasurement("electricity_rates").
			AddTag("rate_type", rateType).
			AddField("value_inc_vat", rate.ValueIncVAT).
			SetTime(rate.ValidFrom)

		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			return err
		}
	}

	return nil
}
