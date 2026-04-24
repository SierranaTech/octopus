package octopus

import (
	"context"
	"net/url"
)

// GridSupplyPoints returns grid supply points for a given postcode.
func (c *Client) GridSupplyPoints(ctx context.Context, postcode string) ([]GridSupplyPoint, error) {
	query := url.Values{"postcode": {postcode}}
	return paginate[GridSupplyPoint](ctx, c, "/industry/grid-supply-points/", query, false)
}
