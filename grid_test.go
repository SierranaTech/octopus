package octopus

import (
	"context"
	"net/http"
	"testing"
)

func TestGridSupplyPoints(t *testing.T) {
	t.Parallel()
	_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/industry/grid-supply-points/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("postcode") != "SW1A1AA" {
			t.Fatalf("unexpected postcode: %s", r.URL.Query().Get("postcode"))
		}
		writeJSON(t, w, Page[GridSupplyPoint]{
			Results: []GridSupplyPoint{{GroupID: "_C"}},
		})
	})
	gsps, err := client.GridSupplyPoints(context.Background(), "SW1A1AA")
	if err != nil {
		t.Fatal(err)
	}
	if len(gsps) != 1 || gsps[0].GroupID != "_C" {
		t.Fatalf("unexpected result: %+v", gsps)
	}
}
