package octopus

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestProducts(t *testing.T) {
	t.Parallel()

	t.Run("ListProducts_NoParams", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/products/" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, Page[Product]{
				Results: []Product{
					{productBase: productBase{Code: "AGILE-24-10-01", DisplayName: "Agile Octopus"}},
					{productBase: productBase{Code: "GO-VAR-22-10-14", DisplayName: "Octopus Go"}},
				},
			})
		})
		products, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(products) != 2 {
			t.Fatalf("expected 2, got %d", len(products))
		}
	})

	t.Run("ListProducts_WithParams", func(t *testing.T) {
		t.Parallel()
		green := true
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("is_green") != "true" {
				t.Fatal("expected is_green=true")
			}
			writeJSON(t, w, Page[Product]{Results: []Product{}})
		})
		_, err := client.ListProducts(context.Background(), &ListProductsParams{IsGreen: &green})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetProduct", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/products/AGILE-24-10-01/" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			writeJSON(t, w, ProductDetail{
				productBase: productBase{Code: "AGILE-24-10-01", DisplayName: "Agile Octopus"},
			})
		})
		detail, err := client.GetProduct(context.Background(), "AGILE-24-10-01", nil)
		if err != nil {
			t.Fatal(err)
		}
		if detail.Code != "AGILE-24-10-01" {
			t.Fatalf("expected AGILE-24-10-01, got %s", detail.Code)
		}
	})

	t.Run("GetProduct_WithTariffActiveAt", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("tariffs_active_at") == "" {
				t.Fatal("expected tariffs_active_at param")
			}
			writeJSON(t, w, ProductDetail{productBase: productBase{Code: "TEST"}})
		})
		at := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		_, err := client.GetProduct(context.Background(), "TEST", &at)
		if err != nil {
			t.Fatal(err)
		}
	})
}
