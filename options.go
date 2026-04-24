package octopus

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Option configures a Client.
type Option func(*Client)

// WithAPIKey sets the API key for authenticated endpoints.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithBaseURL overrides the default API base URL. Useful for testing.
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// ListProductsParams are optional filters for ListProducts.
type ListProductsParams struct {
	Brand       *string
	IsVariable  *bool
	IsBusiness  *bool
	AvailableAt *time.Time
	IsGreen     *bool
	IsPrepay    *bool
}

func (p *ListProductsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Brand != nil {
		v.Set("brand", *p.Brand)
	}
	if p.IsVariable != nil {
		v.Set("is_variable", strconv.FormatBool(*p.IsVariable))
	}
	if p.IsBusiness != nil {
		v.Set("is_business", strconv.FormatBool(*p.IsBusiness))
	}
	if p.AvailableAt != nil {
		v.Set("available_at", p.AvailableAt.Format(time.RFC3339))
	}
	if p.IsGreen != nil {
		v.Set("is_green", strconv.FormatBool(*p.IsGreen))
	}
	if p.IsPrepay != nil {
		v.Set("is_prepay", strconv.FormatBool(*p.IsPrepay))
	}
	return v
}

// TariffParams are optional filters for tariff rate and standing charge endpoints.
type TariffParams struct {
	PeriodFrom *time.Time
	PeriodTo   *time.Time
}

func (p *TariffParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.PeriodFrom != nil {
		v.Set("period_from", p.PeriodFrom.Format(time.RFC3339))
	}
	if p.PeriodTo != nil {
		v.Set("period_to", p.PeriodTo.Format(time.RFC3339))
	}
	return v
}

// ConsumptionParams are optional filters for consumption endpoints.
type ConsumptionParams struct {
	PeriodFrom *time.Time
	PeriodTo   *time.Time
	PageSize   *int
	OrderBy    *string // "period" or "-period"
	GroupBy    *string // "hour", "day", "week", "month", "quarter"
}

func (p *ConsumptionParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.PeriodFrom != nil {
		v.Set("period_from", p.PeriodFrom.Format(time.RFC3339))
	}
	if p.PeriodTo != nil {
		v.Set("period_to", p.PeriodTo.Format(time.RFC3339))
	}
	if p.PageSize != nil {
		v.Set("page_size", strconv.Itoa(*p.PageSize))
	}
	if p.OrderBy != nil {
		v.Set("order_by", *p.OrderBy)
	}
	if p.GroupBy != nil {
		v.Set("group_by", *p.GroupBy)
	}
	return v
}
