package middleware_test

import (
	"github.com/magiconair/properties/assert"
	"github.com/msales/pkg/httpx/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testsData = []struct {
	inUrl  string
	outMap map[string][]string
}{
	{
		inUrl: "/?param[0]=value1&param[1]=value2",
		outMap: map[string][]string{
			"param": {
				"value1",
				"value2",
			},
		},
	},
	{
		inUrl: "/?param[]=value1&param[]=value2",
		outMap: map[string][]string{
			"param": {
				"value1",
				"value2",
			},
		},
	},
	{
		inUrl: "/?param[]=value1&param[]=value2&param2[]=Param1&param2[]=Param2",
		outMap: map[string][]string{
			"param": {
				"value1",
				"value2",
			},
			"region": {
				"Param1",
				"Param2",
			},
		},
	},
	{
		inUrl: "/?param[0]=value1&param[1]=value2&param2[]=Param1&param2[]=Param2",
		outMap: map[string][]string{
			"param": {
				"value1",
				"value2",
			},
			"region": {
				"Param1",
				"Param2",
			},
		},
	},
	{
		inUrl: "/?single=asdf&param[0]=value1&param[1]=value2&param2[]=Param1&param2[]=Param2",
		outMap: map[string][]string{
			"param": {
				"value1",
				"value2",
			},
			"region": {
				"Param1",
				"Param2",
			},
			"single": {
				"asdf",
			},
		},
	},
}

func TestWithQueryNormalizer(t *testing.T) {
	for _, tt := range testsData {
		t.Run(tt.inUrl, func(t *testing.T) {
			h := middleware.WithQueryNormalizer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, r.URL.Query(), tt.outMap)
			}))

			req, _ := http.NewRequest("GET", tt.inUrl, nil)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)
		})
	}
}
