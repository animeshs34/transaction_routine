package api_test

import (
	"github.com/animeshs34/transaction_routine/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Unit test: Recoverer should convert a panic into 500 and not crash the server.
func TestRecoverer_PanicReturns500(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	h := api.Chain(panicHandler, api.Recoverer())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500; got %d", w.Code)
	}
}
