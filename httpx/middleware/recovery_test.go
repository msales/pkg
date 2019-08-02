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
	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/mock"
)

func TestWithRecovery(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "panic", "url", "/", "stack", mock.AnythingOfType("string"))
	s := new(MockStats)
	s.On("Inc", "panic_recovery", int64(1), float32(1.0), mock.Anything).Return(nil).Once()

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(stats.WithStats(log.WithLogger(ctx, logger), s))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)

	logger.AssertExpectations(t)
	s.AssertExpectations(t)
}

func TestWithRecovery_WithoutStack(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}), middleware.WithoutStack())

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "panic", "url", "/")
	s := new(MockStats)
	s.On("Inc", "panic_recovery", int64(1), float32(1.0), mock.Anything).Return(nil).Once()

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(stats.WithStats(log.WithLogger(ctx, logger), s))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)

	logger.AssertExpectations(t)
	s.AssertExpectations(t)
}

func TestWithRecovery_Error(t *testing.T) {
	h := middleware.WithRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("panic"))
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(stats.WithStats(log.WithLogger(context.Background(), log.Null), stats.Null))
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)
}
