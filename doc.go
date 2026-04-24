// Package octopus provides a Go client for the Octopus Energy REST API v1.
//
// Create a client and call methods to interact with the API:
//
//	client := octopus.NewClient(octopus.WithAPIKey("sk_live_..."))
//
//	// List all products (no auth required)
//	products, err := client.ListProducts(ctx, nil)
//
//	// Get electricity consumption (auth required)
//	consumption, err := client.ElectricityConsumption(ctx, mpan, serial, nil)
package octopus
