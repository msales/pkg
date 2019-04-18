package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/v3/httpx/middleware"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/mocks"
	"github.com/stretchr/testify/mock"
)

func TestWithRecovery(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "panic", "stack", mock.AnythingOfType("string"))

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(log.WithLogger(ctx, logger))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)

	logger.AssertExpectations(t)
}

func TestWithRecovery_WithoutStack(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}), middleware.WithoutStack())

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "panic")

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(log.WithLogger(ctx, logger))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)

	logger.AssertExpectations(t)
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
