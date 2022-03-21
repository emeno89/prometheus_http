package middleware

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounters(t *testing.T) {
	h := NewHandler("test_project", func(r *http.Request) string {
		return r.RequestURI
	})

	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Test: %s", r.RequestURI)
	})

	testHandler := h.Middleware(nextHandler)
	testHandler.ServeHTTP(w, r)

	assert.Equal(t, 1.00, testutil.ToFloat64(h.reqCounter))
	assert.Equal(t, 1.00, testutil.ToFloat64(h.respCounter))

	testHandler.ServeHTTP(w, r)

	assert.Equal(t, 2.00, testutil.ToFloat64(h.reqCounter))
	assert.Equal(t, 2.00, testutil.ToFloat64(h.respCounter))
}
