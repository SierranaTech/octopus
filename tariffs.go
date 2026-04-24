package octopus

import (
	"context"
	"fmt"
)

func (c *Client) electricityTariffRates(ctx context.Context, productCode, tariffCode, rateType string, params *TariffParams) ([]TariffCharge, error) {
	path := fmt.Sprintf("/products/%s/electricity-tariffs/%s/%s/", productCode, tariffCode, rateType)
	return paginate[TariffCharge](ctx, c, path, params.values(), false)
}

func (c *Client) gasTariffRates(ctx context.Context, productCode, tariffCode, rateType string, params *TariffParams) ([]TariffCharge, error) {
	path := fmt.Sprintf("/products/%s/gas-tariffs/%s/%s/", productCode, tariffCode, rateType)
	return paginate[TariffCharge](ctx, c, path, params.values(), false)
}

// ElectricityStandingCharges returns standing charges for an electricity tariff.
func (c *Client) ElectricityStandingCharges(ctx context.Context, productCode, tariffCode string, params *TariffParams) ([]TariffCharge, error) {
	return c.electricityTariffRates(ctx, productCode, tariffCode, "standing-charges", params)
}

// ElectricityStandardUnitRates returns standard unit rates for a single-register electricity tariff.
func (c *Client) ElectricityStandardUnitRates(ctx context.Context, productCode, tariffCode string, params *TariffParams) ([]TariffCharge, error) {
	return c.electricityTariffRates(ctx, productCode, tariffCode, "standard-unit-rates", params)
}

// ElectricityDayUnitRates returns day unit rates for a dual-register (Economy 7) electricity tariff.
func (c *Client) ElectricityDayUnitRates(ctx context.Context, productCode, tariffCode string, params *TariffParams) ([]TariffCharge, error) {
	return c.electricityTariffRates(ctx, productCode, tariffCode, "day-unit-rates", params)
}

// ElectricityNightUnitRates returns night unit rates for a dual-register (Economy 7) electricity tariff.
func (c *Client) ElectricityNightUnitRates(ctx context.Context, productCode, tariffCode string, params *TariffParams) ([]TariffCharge, error) {
	return c.electricityTariffRates(ctx, productCode, tariffCode, "night-unit-rates", params)
}

// GasStandingCharges returns standing charges for a gas tariff.
func (c *Client) GasStandingCharges(ctx context.Context, productCode, tariffCode string, params *TariffParams) ([]TariffCharge, error) {
	return c.gasTariffRates(ctx, productCode, tariffCode, "standing-charges", params)
}

// GasStandardUnitRates returns standard unit rates for a gas tariff.
func (c *Client) GasStandardUnitRates(ctx context.Context, productCode, tariffCode string, params *TariffParams) ([]TariffCharge, error) {
	return c.gasTariffRates(ctx, productCode, tariffCode, "standard-unit-rates", params)
}
