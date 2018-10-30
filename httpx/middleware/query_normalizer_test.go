package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"net/url"

	"github.com/msales/pkg/v3/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

var testsData = []struct {
	inUrl  string
	outMap url.Values
}{
	{
		inUrl: "/?param[0]=value1&param[1]=value2",
		outMap: url.Values{
			"param[]": {
				"value1",
				"value2",
			},
		},
	},
	{
		inUrl: "/?param%5B0%5D=value1&param%5B1%5D=value2",
		outMap: url.Values{
			"param[]": {
				"value1",
				"value2",
			},
		},
	},
	{
		inUrl: "/?param[]=value1&param[]=value2",
		outMap: url.Values{
			"param[]": {
				"value1",
				"value2",
			},
		},
	},
	{
		inUrl: "/?param[]=value1&param[]=value2&param2[]=Param1&param2[]=Param2",
		outMap: url.Values{
			"param[]": {
				"value1",
				"value2",
			},
			"param2[]": {
				"Param1",
				"Param2",
			},
		},
	},
	{
		inUrl: "/?param[0]=value1&param[1]=value2&param2[]=Param1&param2[]=Param2",
		outMap: url.Values{
			"param[]": {
				"value1",
				"value2",
			},
			"param2[]": {
				"Param1",
				"Param2",
			},
		},
	},
	{
		inUrl: "/?single=asdf&param[0]=valuea&param[1]=valueb&param2[]=Parama&param2[]=Paramb",
		outMap: url.Values{
			"param[]": {
				"valuea",
				"valueb",
			},
			"param2[]": {
				"Parama",
				"Paramb",
			},
			"single": {
				"asdf",
			},
		},
	},
	{
		inUrl: "/?param[smth][0]=zero&param[smth][0][val]=val",
		outMap: url.Values{
			"param[smth][]": {
				"zero",
			},
			"param[smth][][val]": {
				"val",
			},
		},
	},
	{
		inUrl: "/?param=value",
		outMap: url.Values{
			"param": {
				"value",
			},
		},
	},
}

func TestWithQueryNormalizer(t *testing.T) {
	for _, tt := range testsData {
		t.Run(tt.inUrl, func(t *testing.T) {
			h := middleware.WithQueryNormalizer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				res := map[string][]string(r.URL.Query())
				for k, v := range res {
					mm := tt.outMap[k]
					assert.ElementsMatch(t, v, mm)
				}

			}))

			req, _ := http.NewRequest("GET", tt.inUrl, nil)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)
		})
	}
}

func BenchmarkWithVerySimpleQueryNormalizer(b *testing.B) {

	for n := 0; n < b.N; n++ {
		h := middleware.WithQueryNormalizer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		req, _ := http.NewRequest("GET", "/?param=val", nil)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)
	}
}

func BenchmarkWithQueryBigerNormalizer(b *testing.B) {

	for n := 0; n < b.N; n++ {
		h := middleware.WithQueryNormalizer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		req, _ := http.NewRequest("GET", "/?param[smth][0]=zero&param[smth][0][val][1][asdf][4][asetrhyty][5][ewrytu][7][aerterwert][7]=val", nil)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)
	}
}
