package dmx

import (
	"net/http"
	"testing"
)

var h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//
})

// taken from bmizerany/pat for comparison
func BenchmarkPatternMatchingOneRoute(b *testing.B) {
	mux := New()
	mux.Get("/hello/:name", h)

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		r, err := http.NewRequest("GET", "/hello/blake", nil)
		if err != nil {
			panic(err)
		}
		b.StartTimer()
		mux.ServeHTTP(nil, r)
	}
}

func BenchmarkPatternMatchingMultipleRoutes(b *testing.B) {
	mux := New()
	mux.Get("/a/:name", h)
	mux.Get("/b/:name", h)
	mux.Get("/c/:name", h)
	mux.Get("/d/:name", h)
	mux.Get("/aa/:name", h)
	mux.Get("/bb/:name", h)
	mux.Get("/cc/:name", h)
	mux.Get("/dd/:name", h)
	mux.Get("/aaa/:name", h)
	mux.Get("/bbb/:name", h)
	mux.Get("/ccc/:name", h)
	mux.Get("/ddd/:name", h)
	mux.Get("/aaaa/:name", h)
	mux.Get("/bbbb/:name", h)
	mux.Get("/cccc/:name", h)
	mux.Get("/dddd/:name", h)
	mux.Get("/hello/:name", h)

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		r, err := http.NewRequest("GET", "/hello/blake", nil)
		if err != nil {
			panic(err)
		}
		b.StartTimer()
		mux.ServeHTTP(nil, r)
	}
}

// BenchmarkPatternMatchingOneRoute         2000000               825 ns/op
// BenchmarkPatternMatchingMultipleRoutes   1000000              1360 ns/op
