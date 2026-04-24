package octopus

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// ListProducts returns all available energy products.
func (c *Client) ListProducts(ctx context.Context, params *ListProductsParams) ([]Product, error) {
	var query url.Values
	if params != nil {
		query = params.values()
	}
	return paginate[Product](ctx, c, "/products/", query, false)
}

// GetProduct returns the full details for a single product, including tariffs by region.
// If tariffActiveAt is non-nil, tariff pricing is shown as of that time.
func (c *Client) GetProduct(ctx context.Context, code string, tariffActiveAt *time.Time) (*ProductDetail, error) {
	query := url.Values{}
	if tariffActiveAt != nil {
		query.Set("tariffs_active_at", tariffActiveAt.Format(time.RFC3339))
	}
	var detail ProductDetail
	if err := c.get(ctx, fmt.Sprintf("/products/%s/", code), query, false, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}
