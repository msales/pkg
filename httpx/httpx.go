package httpx

import "net/http"

// HeaderExtractor returns a function that can extract a value from a list of headers.
func HeaderExtractor(headers []string) func(*http.Request) string {
	return func(r *http.Request) string {
		for _, h := range headers {
			v := r.Header.Get(h)
			if v != "" {
				return v
			}
		}

		return ""
	}
}
