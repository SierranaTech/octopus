package octopus

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := NewClient(WithBaseURL(srv.URL), WithAPIKey("test-key"))
	return srv, client
}

func testServerNoAuth(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := NewClient(WithBaseURL(srv.URL))
	return srv, client
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("failed to encode response: %v", err)
	}
}

func TestClient(t *testing.T) {
	t.Parallel()

	t.Run("BasicAuth", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				t.Fatal("expected Basic Auth header")
			}
			if user != "test-key" || pass != "" {
				t.Fatalf("expected user=test-key, pass=empty; got user=%s, pass=%s", user, pass)
			}
			writeJSON(t, w, map[string]any{"number": "A-123"})
		})
		_, err := client.GetAccount(context.Background(), "A-123")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NoAuthHeader_WhenNoAPIKey", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "" {
				t.Fatal("expected no Authorization header")
			}
			writeJSON(t, w, Page[Product]{Results: []Product{}})
		})
		_, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("RequireAuth_ReturnsError", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("should not have made a request")
		})
		_, err := client.GetAccount(context.Background(), "A-123")
		if !errors.Is(err, ErrAuthRequired) {
			t.Fatalf("expected ErrAuthRequired, got %v", err)
		}
	})

	t.Run("NonOKStatus", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"detail":"Not found."}`))
		})
		_, err := client.GetAccount(context.Background(), "A-123")
		if err == nil {
			t.Fatal("expected error")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T", err)
		}
		if apiErr.StatusCode != 404 {
			t.Fatalf("expected status 404, got %d", apiErr.StatusCode)
		}
		if apiErr.Message != "Not found." {
			t.Fatalf("expected parsed message 'Not found.', got %q", apiErr.Message)
		}
	})

	t.Run("IsNotFound", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"detail":"Not found."}`))
		})
		_, err := client.GetAccount(context.Background(), "A-123")
		if !IsNotFound(err) {
			t.Fatal("expected IsNotFound to return true")
		}
	})

	t.Run("IsNotFound_FalseForOtherErrors", func(t *testing.T) {
		t.Parallel()
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"detail":"Server error"}`))
		})
		_, err := client.GetAccount(context.Background(), "A-123")
		if IsNotFound(err) {
			t.Fatal("expected IsNotFound to return false for 500")
		}
	})

	t.Run("SSRFProtection", func(t *testing.T) {
		t.Parallel()
		evil := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request reached evil server — SSRF protection failed")
		}))
		t.Cleanup(evil.Close)

		nextURL := evil.URL + "/steal-creds"
		_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, Page[Product]{
				Results: []Product{{productBase: productBase{Code: "P1"}}},
				Next:    &nextURL,
			})
		})
		_, err := client.ListProducts(context.Background(), nil)
		if err == nil {
			t.Fatal("expected SSRF error, got nil")
		}
	})

	t.Run("AcceptAndUserAgentHeaders", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Accept") != "application/json" {
				t.Fatalf("expected Accept: application/json, got %q", r.Header.Get("Accept"))
			}
			if r.Header.Get("User-Agent") != userAgent {
				t.Fatalf("expected User-Agent: %s, got %q", userAgent, r.Header.Get("User-Agent"))
			}
			writeJSON(t, w, Page[Product]{Results: []Product{}})
		})
		_, err := client.ListProducts(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, Page[Product]{Results: []Product{}})
		})
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := client.ListProducts(ctx, nil)
		if err == nil {
			t.Fatal("expected error from cancelled context")
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		t.Parallel()
		_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{not valid json`))
		})
		_, err := client.ListProducts(context.Background(), nil)
		if err == nil {
			t.Fatal("expected JSON parse error")
		}
	})
}
