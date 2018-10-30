package httpx_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/v3/httpx"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	_, err := httpx.Get(ts.URL)

	assert.NoError(t, err)
}

func TestHead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	_, err := httpx.Head(ts.URL)

	assert.NoError(t, err)
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	body := bytes.NewReader([]byte{})
	_, err := httpx.Post(ts.URL, "plain/text", body)

	assert.NoError(t, err)
}
