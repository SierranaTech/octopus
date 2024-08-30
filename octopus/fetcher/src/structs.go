package octopus

// APIResponse represents the structure of the response from the API
type APIResponse struct {
	Results []TariffRate `json:"results"`
	Next    *string      `json:"next"`
}

// TariffRate represents a single day-unit rate
type TariffRate struct {
	ValueIncVAT float64 `json:"value_inc_vat"`
	StartTime   string  `json:"valid_from"`
	EndTime     string  `json:"valid_to"`
}
