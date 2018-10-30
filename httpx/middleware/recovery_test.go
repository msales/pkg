package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/msales/pkg/v3/httpx/middleware"
	"github.com/msales/pkg/v3/log"
)

func TestWithRecovery(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(log.WithLogger(context.Background(), log.Null))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)
}

func TestWithRecovery_Error(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("panic"))
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(log.WithLogger(context.Background(), log.Null))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)
}
