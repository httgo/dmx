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
	mux.Get("/h/:name", h)
	mux.Get("/he/:name", h)
	mux.Get("/hel/:name", h)
	mux.Get("/hell/:name", h)
	mux.Get("/hellow/:name", h)
	mux.Get("/hellowo/:name", h)
	mux.Get("/hellowor/:name", h)
	mux.Get("/helloworl/:name", h)
	mux.Get("/helloworld/:name", h)
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

// BenchmarkPatternMatchingOneRoute         2000000               934 ns/op
// BenchmarkPatternMatchingMultipleRoutes   1000000              1353 ns/op
