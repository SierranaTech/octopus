# octopus

A Go client for the [Octopus Energy REST API v1](https://developer.octopus.energy/rest/).

Zero dependencies beyond the Go standard library. Covers all 13 public endpoints.

## Install

```sh
go get github.com/sierra1011/octopus@latest
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sierra1011/octopus"
)

func main() {
	// Unauthenticated client — sufficient for product/tariff/grid lookups
	client := octopus.NewClient()

	ctx := context.Background()

	products, err := client.ListProducts(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range products {
		fmt.Printf("%s — %s\n", p.Code, p.DisplayName)
	}
}
```

For endpoints that require authentication (consumption data, account details), pass your API key:

```go
client := octopus.NewClient(octopus.WithAPIKey("sk_live_..."))

consumption, err := client.ElectricityConsumption(ctx, mpan, serialNumber, nil)
```

## Authentication

The Octopus Energy API uses HTTP Basic Auth. Pass your API key via `WithAPIKey` — the client sends it as the Basic Auth username with an empty password. You can find your API key on the [Octopus Energy developer dashboard](https://octopus.energy/dashboard/developer/).

Only three endpoints require authentication: `ElectricityConsumption`, `GasConsumption`, and `GetAccount`. All other endpoints are public. Calling an authenticated endpoint without an API key returns `octopus.ErrAuthRequired`.

## Client options

```go
client := octopus.NewClient(
	octopus.WithAPIKey("sk_live_..."),       // API key for authenticated endpoints
	octopus.WithHTTPClient(customClient),    // custom *http.Client (timeouts, proxies, etc.)
	octopus.WithBaseURL("https://..."),      // override base URL (useful for testing)
)
```

## Endpoints

### Products

```go
// List all available energy products
products, err := client.ListProducts(ctx, nil)

// Filter products
variable := true
products, err := client.ListProducts(ctx, &octopus.ListProductsParams{
	IsVariable: &variable,
	Brand:      strPtr("OCTOPUS_ENERGY"),
})

// Get full details for a single product, including tariffs by region
detail, err := client.GetProduct(ctx, "AGILE-24-10-01", nil)

// Get tariff pricing as of a specific date
at := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
detail, err := client.GetProduct(ctx, "AGILE-24-10-01", &at)
```

### Electricity tariff rates

All tariff methods share the same signature and return `[]TariffCharge`. Pass `nil` for params to get all available data, or filter by time period:

```go
from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
params := &octopus.TariffParams{PeriodFrom: &from, PeriodTo: &to}

// Standing charges
charges, err := client.ElectricityStandingCharges(ctx, productCode, tariffCode, params)

// Standard unit rates (single-register tariffs)
rates, err := client.ElectricityStandardUnitRates(ctx, productCode, tariffCode, params)

// Day unit rates (Economy 7 / dual-register tariffs)
dayRates, err := client.ElectricityDayUnitRates(ctx, productCode, tariffCode, nil)

// Night unit rates (Economy 7 / dual-register tariffs)
nightRates, err := client.ElectricityNightUnitRates(ctx, productCode, tariffCode, nil)
```

### Gas tariff rates

```go
// Standing charges
charges, err := client.GasStandingCharges(ctx, productCode, tariffCode, nil)

// Standard unit rates
rates, err := client.GasStandardUnitRates(ctx, productCode, tariffCode, nil)
```

### Consumption data (authentication required)

```go
client := octopus.NewClient(octopus.WithAPIKey("sk_live_..."))

// Electricity consumption — half-hourly readings
consumption, err := client.ElectricityConsumption(ctx, mpan, serialNumber, nil)

// With filters
from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
pageSize := 48
orderBy := "period"
consumption, err := client.ElectricityConsumption(ctx, mpan, serialNumber, &octopus.ConsumptionParams{
	PeriodFrom: &from,
	PageSize:   &pageSize,   // up to 25,000
	OrderBy:    &orderBy,    // "period" (chronological) or "-period" (reverse)
	GroupBy:    strPtr("day"), // "hour", "day", "week", "month", "quarter"
})

// Gas consumption — same interface
gasData, err := client.GasConsumption(ctx, mprn, serialNumber, nil)
```

### Meter point details

```go
// Look up region and profile class for an MPAN
details, err := client.GetMeterPoint(ctx, "1234567890123")
fmt.Println(details.GSP) // e.g. "_A" for Eastern England
```

### Grid supply points

```go
// Find your region code by postcode
gsps, err := client.GridSupplyPoints(ctx, "SW1A1AA")
fmt.Println(gsps[0].GroupID) // e.g. "_C"
```

### Account details (authentication required)

```go
client := octopus.NewClient(octopus.WithAPIKey("sk_live_..."))

account, err := client.GetAccount(ctx, "A-1234ABCD")
for _, prop := range account.Properties {
	for _, mp := range prop.ElectricityMeterPoints {
		fmt.Printf("MPAN: %s\n", mp.MPAN)
		for _, m := range mp.Meters {
			fmt.Printf("  Meter: %s\n", m.SerialNumber)
		}
	}
}
```

## Pagination

All list endpoints automatically follow pagination links. Results are collected and returned as a single slice. A safety cap of 100 pages prevents runaway loops.

The client validates that pagination `next` URLs share the same origin as the base URL to prevent credential forwarding to untrusted hosts.

## Error handling

```go
rates, err := client.ElectricityStandardUnitRates(ctx, product, tariff, nil)
if err != nil {
	// Check for specific error types
	if err == octopus.ErrAuthRequired {
		// Called an authenticated endpoint without WithAPIKey
	}

	var apiErr *octopus.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API returned %d: %s\n", apiErr.StatusCode, apiErr.Body)
	}
}
```

## Endpoint reference

| Method | Auth | API endpoint |
|--------|------|-------------|
| `ListProducts` | No | `GET /v1/products/` |
| `GetProduct` | No | `GET /v1/products/{code}/` |
| `ElectricityStandingCharges` | No | `GET /v1/products/{code}/electricity-tariffs/{tariff}/standing-charges/` |
| `ElectricityStandardUnitRates` | No | `GET /v1/products/{code}/electricity-tariffs/{tariff}/standard-unit-rates/` |
| `ElectricityDayUnitRates` | No | `GET /v1/products/{code}/electricity-tariffs/{tariff}/day-unit-rates/` |
| `ElectricityNightUnitRates` | No | `GET /v1/products/{code}/electricity-tariffs/{tariff}/night-unit-rates/` |
| `GasStandingCharges` | No | `GET /v1/products/{code}/gas-tariffs/{tariff}/standing-charges/` |
| `GasStandardUnitRates` | No | `GET /v1/products/{code}/gas-tariffs/{tariff}/standard-unit-rates/` |
| `ElectricityConsumption` | Yes | `GET /v1/electricity-meter-points/{mpan}/meters/{serial}/consumption/` |
| `GasConsumption` | Yes | `GET /v1/gas-meter-points/{mprn}/meters/{serial}/consumption/` |
| `GetMeterPoint` | No | `GET /v1/electricity-meter-points/{mpan}/` |
| `GridSupplyPoints` | No | `GET /v1/industry/grid-supply-points/` |
| `GetAccount` | Yes | `GET /v1/accounts/{account}/` |

## Development

Requires [mise](https://mise.jdx.dev/):

```sh
mise install        # install Go 1.26.2 + tools
mise run test       # lint, vet, and run tests with race detection
mise run build      # build all packages
mise run fmt        # format and tidy
```
