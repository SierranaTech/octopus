package octopus

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestTariffs(t *testing.T) {
	t.Parallel()

	t.Run("ElectricityStandingCharges", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "electricity-tariffs") || !strings.Contains(r.URL.Path, "standing-charges") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[TariffCharge]{
				Results: []TariffCharge{{ValueIncVAT: 46.36, ValidFrom: time.Now()}},
			})
		})
		charges, err := client.ElectricityStandingCharges(context.Background(), "PROD", "E-1R-PROD-A", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(charges) != 1 || charges[0].ValueIncVAT != 46.36 {
			t.Fatalf("unexpected result: %+v", charges)
		}
	})

	t.Run("ElectricityStandardUnitRates_WithPeriodParams", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "standard-unit-rates") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.URL.Query().Get("period_from") == "" {
				t.Fatal("expected period_from param")
			}
			writeJSON(t, w, Page[TariffCharge]{Results: []TariffCharge{{ValueIncVAT: 24.5}}})
		})
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		rates, err := client.ElectricityStandardUnitRates(context.Background(), "PROD", "TARIFF", &TariffParams{PeriodFrom: &from})
		if err != nil {
			t.Fatal(err)
		}
		if len(rates) != 1 {
			t.Fatalf("expected 1, got %d", len(rates))
		}
	})

	t.Run("ElectricityDayUnitRates", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "day-unit-rates") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[TariffCharge]{Results: []TariffCharge{{ValueIncVAT: 30.0}}})
		})
		rates, err := client.ElectricityDayUnitRates(context.Background(), "PROD", "TARIFF", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(rates) != 1 {
			t.Fatalf("expected 1, got %d", len(rates))
		}
	})

	t.Run("ElectricityNightUnitRates", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "night-unit-rates") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[TariffCharge]{Results: []TariffCharge{{ValueIncVAT: 10.0}}})
		})
		rates, err := client.ElectricityNightUnitRates(context.Background(), "PROD", "TARIFF", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(rates) != 1 {
			t.Fatalf("expected 1, got %d", len(rates))
		}
	})

	t.Run("GasStandingCharges", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "gas-tariffs") || !strings.Contains(r.URL.Path, "standing-charges") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[TariffCharge]{Results: []TariffCharge{{ValueIncVAT: 30.0}}})
		})
		charges, err := client.GasStandingCharges(context.Background(), "PROD", "G-1R-PROD-A", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(charges) != 1 {
			t.Fatalf("expected 1, got %d", len(charges))
		}
	})

	t.Run("GasStandardUnitRates", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "gas-tariffs") || !strings.Contains(r.URL.Path, "standard-unit-rates") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[TariffCharge]{Results: []TariffCharge{{ValueIncVAT: 7.5}}})
		})
		rates, err := client.GasStandardUnitRates(context.Background(), "PROD", "G-1R-PROD-A", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(rates) != 1 {
			t.Fatalf("expected 1, got %d", len(rates))
		}
	})

	t.Run("NullValidTo", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"count":1,"next":null,"previous":null,"results":[{"value_exc_vat":20.0,"value_inc_vat":24.0,"valid_from":"2024-01-01T00:00:00Z","valid_to":null}]}`))
		})
		rates, err := client.ElectricityStandardUnitRates(context.Background(), "PROD", "TARIFF", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(rates) != 1 {
			t.Fatalf("expected 1, got %d", len(rates))
		}
		if rates[0].ValidTo != nil {
			t.Fatal("expected ValidTo to be nil for current rate")
		}
		if rates[0].ValueExcVAT != 20.0 {
			t.Fatalf("expected 20.0, got %f", rates[0].ValueExcVAT)
		}
	})
}
