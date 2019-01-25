package httpx

import (
	"io"
	"net"
	"net/http"
	"time"
)

// DefaultTransport is the default implementation of Transport and is
// used by DefaultClient.
var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 90 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout: 3 * time.Second,
}

// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &http.Client{
	Transport: DefaultTransport,
}

// Get issues a GET to the specified URL.
func Get(url string) (resp *http.Response, err error) {
	return DefaultClient.Get(url)
}

// Head issues a HEAD to the specified URL.
func Head(url string) (resp *http.Response, err error) {
	return DefaultClient.Head(url)
}

// Post issues a POST to the specified URL.
func Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return DefaultClient.Post(url, contentType, body)
}
