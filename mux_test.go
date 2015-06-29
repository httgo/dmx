package dmx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/nowk/assert.v2"
)

func hFunc(str string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(str))
	}
}

func mFunc(a, b string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(a))
			next.ServeHTTP(w, req)
			w.Write([]byte(b))
		})
	}
}

func send(
	t *testing.T,
	m *Mux,
	meth, path string,
	b io.Reader) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()

	req, err := http.NewRequest(meth, path, b)
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(w, req)

	return w
}

var auth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("auth") == "true" {
			next.ServeHTTP(w, req)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}

var methodOverride = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("_method") == "PUT" {
			req.Method = "PUT"
		}

		next.ServeHTTP(w, req)
	})
}

var stop = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("<stop>"))
	})
}

func TestBasicRouting(t *testing.T) {
	mux := New()
	mux.GetFunc("/a", hFunc("a"))
	mux.GetFunc("/b", hFunc("b"))
	mux.GetFunc("/c", hFunc("c"))

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/a", "a"},
		{"GET", "/b", "b"},
		{"GET", "/c", "c"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, mux, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String())
	}
}

func TestBasicRoutingThroughMountedMux(t *testing.T) {
	c := New()
	c.GetFunc("/c", hFunc("c"))

	b := New()
	b.GetFunc("/b", hFunc("b"))
	b.Mount(c)

	a := New()
	a.GetFunc("/a", hFunc("a"))
	a.Mount(b)

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/a", "a"},
		{"GET", "/b", "b"},
		{"GET", "/c", "c"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, a, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String())
	}
}

func TestNotFoundAlwaysBubblesUpToTheMainMuxsElse(t *testing.T) {
	t.Skip("TODO")
}

func TestMiddlewaresAreExtendedDownTheMountChain(t *testing.T) {
	c := New()
	c.Use(mFunc("c", ""))
	c.GetFunc("/d", hFunc("d"))

	b := New()
	b.Use(mFunc("b", ""))
	b.Mount(c)

	a := New()
	a.Use(mFunc("a", ""))
	a.Mount(b)

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/d", "abcd"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, a, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String())
	}
}

func TestMiddlewaresWrapEachOTher(t *testing.T) {
	c := New()
	c.Use(mFunc("c", "c"))
	c.GetFunc("/d", hFunc("d"))

	b := New()
	b.Use(mFunc("b", "b"))
	b.Mount(c)

	a := New()
	a.Use(mFunc("a", "a"))
	a.Mount(b)

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/d", "abcdcba"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, a, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String())
	}
}

func TestMountingCopiesTheIncomingMuxLeavingOriginalUntouched(t *testing.T) {
	t.Skip("TODO")
}

func TestMiddlewaresAreNextable(t *testing.T) {
	mux := New()
	mux.Use(mFunc("a", ""))
	mux.Use(mFunc("b", ""))
	mux.GetFunc("/c", hFunc("c"))

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/c", "abc"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, mux, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String())
	}
}

func TestMiddlewareStopsRequestWhenNoNext(t *testing.T) {
	mux := New()
	mux.Use(mFunc("a", ""))
	mux.Use(stop)
	mux.Use(mFunc("b", ""))
	mux.GetFunc("/c", hFunc("c"))

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/c", "a<stop>"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, mux, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String())
	}
}

func TestNotFoundDoesNotExecuteMiddlewares(t *testing.T) {
	mux := New()
	mux.Use(mFunc("a", "a"))
	mux.Use(stop)

	w := send(t, mux, "GET", "/does/not/exist", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "404 page not found\n", w.Body.String())
}

func TestMustIsCalledBeforeMatch(t *testing.T) {
	mux := New()
	mux.Must(methodOverride)
	mux.PutFunc("/a", hFunc("a"))

	w := send(t, mux, "POST", "/a?_method=PUT", nil)
	assert.Equal(t, "a", w.Body.String())
}

func TestMustIsCalledOnNotFound(t *testing.T) {
	mux := New()
	mux.Must(mFunc("-- ", " --"))

	w := send(t, mux, "GET", "/not/found", nil)
	assert.Equal(t, "-- 404 page not found\n --", w.Body.String())
}

func TestWalkingIntoMultipleMuxRespectsParentMiddlewares(t *testing.T) {
	b := New()
	b.GetFunc("/b", hFunc("b"))

	c := New()
	c.GetFunc("/c", hFunc("c"))

	a := New()
	a.Use(auth)
	a.GetFunc("/a", hFunc("a"))
	a.Mount(b)
	a.Mount(c)

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/c", "not authorized"},
		{"GET", "/b", "not authorized"},
		{"GET", "/a", "not authorized"},
		{"GET", "/c?auth=true", "c"},
		{"GET", "/b?auth=true", "b"},
		{"GET", "/a?auth=true", "a"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, a, v.m, v.p, nil)

		assert.Equal(t, v.b, w.Body.String(), v.m, v.p)
	}
}

func TestMultiAuthsThroughNestedMounts(t *testing.T) {
	b := New()
	b.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Query().Get("auth2") == "true" {
				next.ServeHTTP(w, req)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized - 2"))
			}
		})
	})
	b.GetFunc("/b", hFunc("b"))

	c := New()
	c.GetFunc("/c", hFunc("c"))

	a := New()
	a.Use(auth)
	a.GetFunc("/a", hFunc("a"))
	a.Mount(b)
	a.Mount(c)

	for _, v := range []struct {
		m, p, b string
	}{
		{"GET", "/b?auth=true", "not authorized - 2"},
		{"GET", "/b?auth=true&auth2=true", "b"},
		{"GET", "/b?auth2=true", "not authorized"},
		{"GET", "/c?auth=true", "c"},
		{"GET", "/c?auth=true&auth2=true", "c"},
		{"GET", "/c?auth2=true", "not authorized"},
		{"GET", "/a?auth=true", "a"},
		{"GET", "/a?auth=true&auth2=true", "a"},
		{"GET", "/a?auth2=true", "not authorized"},

		{"GET", "/not/found", "404 page not found\n"},
	} {
		w := send(t, a, v.m, v.p, nil)
		assert.Equal(t, v.b, w.Body.String(), v.m, v.p)
	}
}
