package octopus

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestConsumption(t *testing.T) {
	t.Parallel()

	t.Run("ElectricityConsumption", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "electricity-meter-points") || !strings.Contains(r.URL.Path, "consumption") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			user, _, ok := r.BasicAuth()
			if !ok || user != "test-key" {
				t.Fatal("expected Basic Auth")
			}
			writeJSON(t, w, Page[Consumption]{
				Results: []Consumption{
					{Consumption: 0.523, IntervalStart: time.Now(), IntervalEnd: time.Now()},
				},
			})
		})
		data, err := client.ElectricityConsumption(context.Background(), "1234567890123", "19L1234567", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) != 1 || data[0].Consumption != 0.523 {
			t.Fatalf("unexpected result: %+v", data)
		}
	})

	t.Run("ElectricityConsumption_RequiresAuth", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("should not make request")
		})
		_, err := client.ElectricityConsumption(context.Background(), "MPAN", "SERIAL", nil)
		if !errors.Is(err, ErrAuthRequired) {
			t.Fatalf("expected ErrAuthRequired, got %v", err)
		}
	})

	t.Run("GasConsumption", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "gas-meter-points") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[Consumption]{
				Results: []Consumption{{Consumption: 1.234}},
			})
		})
		data, err := client.GasConsumption(context.Background(), "1234567890", "G4A12345", nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) != 1 {
			t.Fatalf("expected 1, got %d", len(data))
		}
	})

	t.Run("GasConsumption_RequiresAuth", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("should not make request")
		})
		_, err := client.GasConsumption(context.Background(), "MPRN", "SERIAL", nil)
		if !errors.Is(err, ErrAuthRequired) {
			t.Fatalf("expected ErrAuthRequired, got %v", err)
		}
	})

	t.Run("ElectricityConsumptionPages", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, Page[Consumption]{
				Results: []Consumption{{Consumption: 0.5}, {Consumption: 0.6}},
			})
		})
		var total int
		for page, err := range client.ElectricityConsumptionPages(context.Background(), "MPAN", "SERIAL", nil) {
			if err != nil {
				t.Fatal(err)
			}
			total += len(page)
		}
		if total != 2 {
			t.Fatalf("expected 2, got %d", total)
		}
	})

	t.Run("WithAllParams", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("period_from") == "" {
				t.Fatal("expected period_from")
			}
			if q.Get("page_size") != "48" {
				t.Fatalf("expected page_size=48, got %s", q.Get("page_size"))
			}
			if q.Get("order_by") != "period" {
				t.Fatal("expected order_by=period")
			}
			if q.Get("group_by") != "day" {
				t.Fatalf("expected group_by=day, got %s", q.Get("group_by"))
			}
			writeJSON(t, w, Page[Consumption]{Results: []Consumption{}})
		})
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		pageSize := 48
		orderBy := "period"
		groupBy := "day"
		_, err := client.ElectricityConsumption(context.Background(), "MPAN", "SERIAL", &ConsumptionParams{
			PeriodFrom: &from,
			PageSize:   &pageSize,
			OrderBy:    &orderBy,
			GroupBy:    &groupBy,
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}
