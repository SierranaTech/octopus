package octopus

import (
	"context"
	"net/http"
	"testing"
)

func TestGetMeterPoint(t *testing.T) {
	t.Parallel()
	_, client := testServerNoAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/electricity-meter-points/1234567890123/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(t, w, MeterPointDetails{
			GSP:          "_A",
			MPAN:         "1234567890123",
			ProfileClass: 1,
		})
	})
	details, err := client.GetMeterPoint(context.Background(), "1234567890123")
	if err != nil {
		t.Fatal(err)
	}
	if details.GSP != "_A" {
		t.Fatalf("expected _A, got %s", details.GSP)
	}
	if details.ProfileClass != 1 {
		t.Fatalf("expected profile class 1, got %d", details.ProfileClass)
	}
}
