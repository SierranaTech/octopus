package octopus

import "time"

// Page is a generic paginated response from the Octopus Energy API.
type Page[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

// productBase contains fields shared between Product and ProductDetail.
type productBase struct {
	Code          string     `json:"code"`
	FullName      string     `json:"full_name"`
	DisplayName   string     `json:"display_name"`
	Description   string     `json:"description"`
	IsVariable    bool       `json:"is_variable"`
	IsBusiness    bool       `json:"is_business"`
	IsGreen       bool       `json:"is_green"`
	IsPrepay      bool       `json:"is_prepay"`
	IsRestricted  bool       `json:"is_restricted"`
	Brand         string     `json:"brand"`
	Term          *int       `json:"term"`
	AvailableFrom *time.Time `json:"available_from"`
	AvailableTo   *time.Time `json:"available_to"`
	Direction     string     `json:"direction"`
}

// Product represents an energy product from the products listing.
type Product struct {
	productBase
	LinksRel []Link `json:"links"`
}

// Link is an API hyperlink reference.
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

// ProductDetail is the full detail for a single product, including tariffs by region.
type ProductDetail struct {
	productBase

	SingleRegisterElectricityTariffs map[string]RegionTariff `json:"single_register_electricity_tariffs"`
	DualRegisterElectricityTariffs   map[string]RegionTariff `json:"dual_register_electricity_tariffs"`
	SingleRegisterGasTariffs         map[string]RegionTariff `json:"single_register_gas_tariffs"`
}

// RegionTariff holds the tariff details for a specific region.
type RegionTariff struct {
	DirectDebitMonthly   *TariffDetail `json:"direct_debit_monthly"`
	DirectDebitQuarterly *TariffDetail `json:"direct_debit_quarterly"`
	Varying              *TariffDetail `json:"varying"`
}

// TariffDetail contains rate information for a specific payment method.
type TariffDetail struct {
	Code                   string  `json:"code"`
	StandingChargeExcVAT   float64 `json:"standing_charge_exc_vat"`
	StandingChargeIncVAT   float64 `json:"standing_charge_inc_vat"`
	OnlineDiscountExcVAT   float64 `json:"online_discount_exc_vat"`
	OnlineDiscountIncVAT   float64 `json:"online_discount_inc_vat"`
	DualFuelDiscountExcVAT float64 `json:"dual_fuel_discount_exc_vat"`
	DualFuelDiscountIncVAT float64 `json:"dual_fuel_discount_inc_vat"`
	ExitFeesExcVAT         float64 `json:"exit_fees_exc_vat"`
	ExitFeesIncVAT         float64 `json:"exit_fees_inc_vat"`
	ExitFeesType           string  `json:"exit_fees_type"`
	LinksRel               []Link  `json:"links"`
}

// TariffCharge represents a single rate or standing charge entry.
type TariffCharge struct {
	ValueExcVAT float64    `json:"value_exc_vat"`
	ValueIncVAT float64    `json:"value_inc_vat"`
	ValidFrom   time.Time  `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
}

// Consumption represents a single consumption reading.
type Consumption struct {
	Consumption   float64   `json:"consumption"`
	IntervalStart time.Time `json:"interval_start"`
	IntervalEnd   time.Time `json:"interval_end"`
}

// MeterPointDetails holds information about an electricity meter point.
type MeterPointDetails struct {
	GSP          string `json:"gsp"`
	MPAN         string `json:"mpan"`
	ProfileClass int    `json:"profile_class"`
}

// GridSupplyPoint represents a grid supply point returned by postcode lookup.
type GridSupplyPoint struct {
	GroupID string `json:"group_id"`
}

// Account represents an Octopus Energy customer account.
type Account struct {
	Number     string     `json:"number"`
	Properties []Property `json:"properties"`
}

// Property represents a property on an account.
type Property struct {
	ID                     int                     `json:"id"`
	MovedInAt              *time.Time              `json:"moved_in_at"`
	MovedOutAt             *time.Time              `json:"moved_out_at"`
	AddressLine1           string                  `json:"address_line_1"`
	AddressLine2           string                  `json:"address_line_2"`
	AddressLine3           string                  `json:"address_line_3"`
	Town                   string                  `json:"town"`
	County                 string                  `json:"county"`
	Postcode               string                  `json:"postcode"`
	ElectricityMeterPoints []ElectricityMeterPoint `json:"electricity_meter_points"`
	GasMeterPoints         []GasMeterPoint         `json:"gas_meter_points"`
}

// ElectricityMeterPoint represents an electricity meter point on a property.
type ElectricityMeterPoint struct {
	MPAN                string            `json:"mpan"`
	ProfileClass        int               `json:"profile_class"`
	ConsumptionStandard int               `json:"consumption_standard"`
	Meters              []Meter           `json:"meters"`
	Agreements          []TariffAgreement `json:"agreements"`
	IsExport            bool              `json:"is_export"`
}

// GasMeterPoint represents a gas meter point on a property.
type GasMeterPoint struct {
	MPRN                string            `json:"mprn"`
	ConsumptionStandard int               `json:"consumption_standard"`
	Meters              []Meter           `json:"meters"`
	Agreements          []TariffAgreement `json:"agreements"`
}

// Meter represents a physical meter.
type Meter struct {
	SerialNumber string     `json:"serial_number"`
	Registers    []Register `json:"registers"`
}

// Register represents a meter register.
type Register struct {
	Identifier           string `json:"identifier"`
	Rate                 string `json:"rate"`
	IsSettlementRegister bool   `json:"is_settlement_register"`
}

// TariffAgreement represents a tariff agreement period.
type TariffAgreement struct {
	TariffCode string     `json:"tariff_code"`
	ValidFrom  time.Time  `json:"valid_from"`
	ValidTo    *time.Time `json:"valid_to"`
}
