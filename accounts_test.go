package octopus

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestAccounts(t *testing.T) {
	t.Parallel()

	t.Run("GetAccount", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/accounts/A-1234ABCD/" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			user, _, ok := r.BasicAuth()
			if !ok || user != "test-key" {
				t.Fatal("expected Basic Auth with test-key")
			}
			writeJSON(t, w, Account{
				Number: "A-1234ABCD",
				Properties: []Property{
					{
						ID:       1,
						Postcode: "SW1A 1AA",
						ElectricityMeterPoints: []ElectricityMeterPoint{
							{MPAN: "1234567890123", Meters: []Meter{{SerialNumber: "19L1234567"}}},
						},
					},
				},
			})
		})
		account, err := client.GetAccount(context.Background(), "A-1234ABCD")
		if err != nil {
			t.Fatal(err)
		}
		if account.Number != "A-1234ABCD" {
			t.Fatalf("expected A-1234ABCD, got %s", account.Number)
		}
		if len(account.Properties) != 1 {
			t.Fatalf("expected 1 property, got %d", len(account.Properties))
		}
		if account.Properties[0].ElectricityMeterPoints[0].MPAN != "1234567890123" {
			t.Fatal("unexpected MPAN")
		}
	})

	t.Run("GetAccount_RequiresAuth", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("should not make request")
		})
		_, err := client.GetAccount(context.Background(), "A-123")
		if !errors.Is(err, ErrAuthRequired) {
			t.Fatalf("expected ErrAuthRequired, got %v", err)
		}
	})
}
