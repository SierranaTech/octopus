package octopus

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

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
