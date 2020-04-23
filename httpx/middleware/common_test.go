package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/v4/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestWithCommon(t *testing.T) {
	var called bool

	h := middleware.WithCommon(context.Background(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)

	assert.True(t, called)
}
