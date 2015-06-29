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
	mux.GetFunc("/hello/blake", h)

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
	mux.GetFunc("/hello/:name", h)

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
	a.GetFunc("/hello/:name", h)

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
	mux.GetFunc("/h/:name", h)
	mux.GetFunc("/he/:name", h)
	mux.GetFunc("/hel/:name", h)
	mux.GetFunc("/hell/:name", h)
	mux.GetFunc("/hellow/:name", h)
	mux.GetFunc("/hellowo/:name", h)
	mux.GetFunc("/hellowor/:name", h)
	mux.GetFunc("/helloworl/:name", h)
	mux.GetFunc("/helloworld/:name", h)
	mux.GetFunc("/hello/:name", h)

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
