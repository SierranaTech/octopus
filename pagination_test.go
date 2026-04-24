package octopus

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestPagination(t *testing.T) {
	t.Parallel()

	t.Run("MultiplePages", func(t *testing.T) {
		t.Parallel()
		var serverURL atomic.Value
		var page atomic.Int32

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := page.Add(1)
			base := serverURL.Load().(string)
			switch p {
			case 1:
				next := base + "/products/?page=2"
				writeJSON(t, w, Page[Product]{
					Count:   3,
					Results: []Product{{productBase: productBase{Code: "P1"}}, {productBase: productBase{Code: "P2"}}},
					Next:    &next,
				})
			case 2:
				writeJSON(t, w, Page[Product]{
					Count:   3,
					Results: []Product{{productBase: productBase{Code: "P3"}}},
					Next:    nil,
				})
			default:
				t.Error("unexpected page request")
			}
		}))
		t.Cleanup(srv.Close)
		serverURL.Store(srv.URL)

		client := NewClient(WithBaseURL(srv.URL))
		products, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(products) != 3 {
			t.Fatalf("expected 3 products, got %d", len(products))
		}
		if products[2].Code != "P3" {
			t.Fatalf("expected P3, got %s", products[2].Code)
		}
	})

	t.Run("SinglePage", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, Page[Product]{
				Count:   1,
				Results: []Product{{productBase: productBase{Code: "P1"}}},
				Next:    nil,
			})
		})
		products, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(products) != 1 {
			t.Fatalf("expected 1 product, got %d", len(products))
		}
	})

	t.Run("EmptyResults", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, Page[Product]{
				Count:   0,
				Results: []Product{},
				Next:    nil,
			})
		})
		products, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(products) != 0 {
			t.Fatalf("expected 0 products, got %d", len(products))
		}
	})

	t.Run("EmptyNextString", func(t *testing.T) {
		t.Parallel()
		empty := ""
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, Page[Product]{
				Count:   1,
				Results: []Product{{productBase: productBase{Code: "P1"}}},
				Next:    &empty,
			})
		})
		products, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(products) != 1 {
			t.Fatalf("expected 1 product, got %d", len(products))
		}
	})

	t.Run("PageLimitExceeded", func(t *testing.T) {
		t.Parallel()
		var serverURL atomic.Value

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			base := serverURL.Load().(string)
			next := base + "/products/?page=next"
			writeJSON(t, w, Page[Product]{
				Count:   99999,
				Results: []Product{{productBase: productBase{Code: "P"}}},
				Next:    &next,
			})
		}))
		t.Cleanup(srv.Close)
		serverURL.Store(srv.URL)

		client := NewClient(WithBaseURL(srv.URL))
		products, err := client.ListProducts(context.Background(), nil)
		if !errors.Is(err, ErrPageLimitExceeded) {
			t.Fatalf("expected ErrPageLimitExceeded, got %v", err)
		}
		if len(products) == 0 {
			t.Fatal("expected partial results before page limit")
		}
	})
}
