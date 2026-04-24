package octopus

import (
	"context"
	"fmt"
	"iter"
)

// ElectricityConsumption returns half-hourly consumption data for an electricity meter.
// Requires an API key.
func (c *Client) ElectricityConsumption(ctx context.Context, mpan, serialNumber string, params *ConsumptionParams) ([]Consumption, error) {
	if err := c.requireAuth(); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/electricity-meter-points/%s/meters/%s/consumption/", mpan, serialNumber)
	return paginate[Consumption](ctx, c, path, params.values(), true)
}

// ElectricityConsumptionPages returns an iterator that yields one page of consumption data
// at a time, allowing streaming processing without loading all pages into memory.
// Requires an API key.
func (c *Client) ElectricityConsumptionPages(ctx context.Context, mpan, serialNumber string, params *ConsumptionParams) iter.Seq2[[]Consumption, error] {
	if err := c.requireAuth(); err != nil {
		return func(yield func([]Consumption, error) bool) {
			yield(nil, err)
		}
	}
	path := fmt.Sprintf("/electricity-meter-points/%s/meters/%s/consumption/", mpan, serialNumber)
	return pages[Consumption](ctx, c, path, params.values(), true)
}

// GasConsumption returns half-hourly consumption data for a gas meter.
// Requires an API key.
func (c *Client) GasConsumption(ctx context.Context, mprn, serialNumber string, params *ConsumptionParams) ([]Consumption, error) {
	if err := c.requireAuth(); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/gas-meter-points/%s/meters/%s/consumption/", mprn, serialNumber)
	return paginate[Consumption](ctx, c, path, params.values(), true)
}

// GasConsumptionPages returns an iterator that yields one page of gas consumption data
// at a time, allowing streaming processing without loading all pages into memory.
// Requires an API key.
func (c *Client) GasConsumptionPages(ctx context.Context, mprn, serialNumber string, params *ConsumptionParams) iter.Seq2[[]Consumption, error] {
	if err := c.requireAuth(); err != nil {
		return func(yield func([]Consumption, error) bool) {
			yield(nil, err)
		}
	}
	path := fmt.Sprintf("/gas-meter-points/%s/meters/%s/consumption/", mprn, serialNumber)
	return pages[Consumption](ctx, c, path, params.values(), true)
}
