package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msales/pkg/v4/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")

	h := middleware.WithContext(ctx, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test", r.Context().Value("test"))
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)
}

func TestWithContext_Value(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")
	reqCtx := context.WithValue(context.Background(), "req", "req")

	h := middleware.WithContext(ctx, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test", r.Context().Value("test"))
		assert.Equal(t, "req", r.Context().Value("req"))
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req.WithContext(reqCtx))
}

func TestWithContext_WithDeadline(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")

	deadline, err := time.Parse(time.RFC822Z, "11 Nov 19 17:32 -0700")
	if err != nil {
		panic(err)
	}

	reqCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	h := middleware.WithContext(ctx, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test", r.Context().Value("test"))

		reqDeadline, ok := r.Context().Deadline()
		assert.True(t, ok)
		assert.Equal(t, deadline, reqDeadline)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req.WithContext(reqCtx))
}

func TestWithContext_WithCancel(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")
	reqCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := middleware.WithContext(ctx, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test", r.Context().Value("test"))

		reqDone := r.Context().Done()
		assert.NotNil(t, reqDone)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req.WithContext(reqCtx))
}
