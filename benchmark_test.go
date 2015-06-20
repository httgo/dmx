package dmx

import (
	"net/http"
	"testing"
)

var h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//
})

func BenchmarkMatchExactRoute(b *testing.B) {
	mux := New()
	mux.GETFunc("/hello/blake", h)

	b.ReportAllocs()
	b.ResetTimer()

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

func BenchmarkMatchOneRouteWithOneParam(b *testing.B) {
	mux := New()
	mux.GETFunc("/hello/:name", h)

	b.ReportAllocs()
	b.ResetTimer()

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

func BenchmarkMatchOneRouteWithOneParamMountedToOneMiddleware(b *testing.B) {
	a := New()
	a.GETFunc("/hello/:name", h)

	mux := New()
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	})
	mux.Mount(a)

	b.ReportAllocs()
	b.ResetTimer()

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

func BenchmarkMatchRouteWithOneParamAtTheEndOfaListOfSimilarPaths(
	b *testing.B) {

	mux := New()
	mux.GETFunc("/h/:name", h)
	mux.GETFunc("/he/:name", h)
	mux.GETFunc("/hel/:name", h)
	mux.GETFunc("/hell/:name", h)
	mux.GETFunc("/hellow/:name", h)
	mux.GETFunc("/hellowo/:name", h)
	mux.GETFunc("/hellowor/:name", h)
	mux.GETFunc("/helloworl/:name", h)
	mux.GETFunc("/helloworld/:name", h)
	mux.GETFunc("/hello/:name", h)

	b.ReportAllocs()
	b.ResetTimer()

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

// Chrome Pixel 2015 (LS) - Secure shell - go1.4.2 linux/amd64
// PASS
// BenchmarkMatchExactRoute                                        20000000               102 ns/op               0 B/op          0 allocs/op
// BenchmarkMatchOneRouteWithOneParam                               2000000               816 ns/op              48 B/op          2 allocs/op
// BenchmarkMatchOneRouteWithOneParamMountedToOneMiddleware         2000000               870 ns/op              48 B/op          2 allocs/op
// BenchmarkMatchRouteWithOneParamAtTheEndOfaListOfSimilarPaths     1000000              1053 ns/op              48 B/op          2 allocs/op
// ok      github.com/httgo/dmxr   115.220s
