package octopus

import (
	"context"
	"fmt"
)

// GetAccount returns details for a customer account, including properties and meters.
// Requires an API key.
func (c *Client) GetAccount(ctx context.Context, accountNumber string) (*Account, error) {
	if err := c.requireAuth(); err != nil {
		return nil, err
	}
	var account Account
	if err := c.get(ctx, fmt.Sprintf("/accounts/%s/", accountNumber), nil, true, &account); err != nil {
		return nil, err
	}
	return &account, nil
}
