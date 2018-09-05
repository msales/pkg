package middleware

import (
	"net/http"

	"net/url"
	"regexp"
)

// WithQueryNormalizer fixes wrong php query string handling as array
func WithQueryNormalizer(h http.Handler) http.Handler {
	rxp, e := regexp.Compile("\\[[0-9]+\\]")
	if e != nil {
		panic(e)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		normalizedQuery := make(url.Values, 0)

		for key, qVal := range q {
			newKey, values := getNormalizedValue(key, qVal, rxp)

			for _, v := range values {
				normalizedQuery.Add(newKey, v)
			}
		}

		r.URL.RawQuery = normalizedQuery.Encode()
	})
}

func getNormalizedValue(key string, qVal []string, rxp *regexp.Regexp) (string, []string) {
	isNastyArray := rxp.Match([]byte(key))

	if isNastyArray == false {
		return key, qVal
	}

	return rxp.ReplaceAllString(key, ""), qVal
}
