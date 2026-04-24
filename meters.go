package octopus

import (
	"context"
	"fmt"
)

// GetMeterPoint returns details for an electricity meter point (MPAN), including its region.
func (c *Client) GetMeterPoint(ctx context.Context, mpan string) (*MeterPointDetails, error) {
	var details MeterPointDetails
	if err := c.get(ctx, fmt.Sprintf("/electricity-meter-points/%s/", mpan), nil, false, &details); err != nil {
		return nil, err
	}
	return &details, nil
}
