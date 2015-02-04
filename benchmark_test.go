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
			b.Fatal(err)
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
			b.Fatal(err)
		}
		b.StartTimer()
		mux.ServeHTTP(nil, r)
	}
}

func BenchmarkPatternMatchingOneRouteWithFormat(b *testing.B) {
	mux := New()
	mux.Get("/hello/:name.:format", h)

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		r, err := http.NewRequest("GET", "/hello/blake.html", nil)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		mux.ServeHTTP(nil, r)
	}
}

func BenchmarkPatternMatchingMultipleRoutesWithFormat(b *testing.B) {
	mux := New()
	mux.Get("/h/:name", h)
	mux.Get("/he/:name.:format", h)
	mux.Get("/hel/:name", h)
	mux.Get("/hell/:name.:format", h)
	mux.Get("/hellow/:name", h)
	mux.Get("/hellowo/:name.:format", h)
	mux.Get("/hellowor/:name", h)
	mux.Get("/helloworl/:name.:format", h)
	mux.Get("/helloworld/:name", h)
	mux.Get("/hello/:name.:format", h)

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		r, err := http.NewRequest("GET", "/hello/blake.html", nil)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		mux.ServeHTTP(nil, r)
	}
}

// BenchmarkPatternMatchingOneRoute                         1000000              1085 ns/op
// BenchmarkPatternMatchingMultipleRoutes                   1000000              1621 ns/op
// BenchmarkPatternMatchingOneRouteWithFormat               1000000              1342 ns/op
// BenchmarkPatternMatchingMultipleRoutesWithFormat         1000000              2581 ns/op
