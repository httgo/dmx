package mountn

import (
	"fmt"
	"github.com/nowk/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func hfunc(s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(s))
	})
}

func TestMethodPatDuplicationPanics(t *testing.T) {
	for _, v := range []string{
		"",
		"/posts",
	} {
		s := v
		if s == "" {
			s = "/"
		}

		{
			mux := New()
			assert.Panic(t, fmt.Errorf("error: mux: POST %s is already defined", s),
				func() {
					mux.Add(v, hfunc(""), "POST", "POST")
				})
			assert.Equal(t, 1, mux.Len())
		}
		{
			mux := New()
			assert.Panic(t, fmt.Errorf("error: mux: POST %s is already defined", s),
				func() {
					mux.Add(v+"/", hfunc(""), "POST")
					mux.Add(v, hfunc(""), "POST")
				})
			assert.Equal(t, 1, mux.Len())
		}
		{
			mux := New()
			assert.Panic(t, fmt.Errorf("error: mux: POST %s is already defined", s),
				func() {
					mux.Add(v, hfunc(""), "POST")
					mux.Add(v+"/", hfunc(""), "POST")
				})
			assert.Equal(t, 1, mux.Len())
		}
	}
}

func TestDispatchesToMatchingResource(t *testing.T) {
	mux := New()
	mux.Add("/", hfunc("GET /"), "GET")
	mux.Add("/posts/:id", hfunc("POST /posts/:id"), "POST", "PUT")
	mux.Add("/posts/:id", hfunc("GET /posts/:id"), "GET")

	for _, v := range []struct {
		m, u, b string
	}{
		{"GET", "/", "GET /"},
		{"GET", "/posts", "Not Found\n"},
		{"GET", "/posts/", "Not Found\n"},
		{"GET", "/posts/123", "GET /posts/:id"},
		{"GET", "/posts/123/", "GET /posts/:id"},
		{"POST", "/posts/123", "POST /posts/:id"},
		{"POST", "/posts/123/", "POST /posts/:id"},
		{"PUT", "/posts/123", "POST /posts/:id"},
		{"PUT", "/posts/123/", "POST /posts/:id"},
	} {
		u := fmt.Sprintf("http://www.com%s", v.u)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(v.m, u, nil)
		if err != nil {
			t.Fatal(err)
		}

		mux.ServeHTTP(w, req)
		assert.Equal(t, v.b, w.Body.String(), v.m, ":", v.u, ":", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New()
	mux.Add("/posts/:id", hfunc(""), "POST", "GET", "DELETE")

	w := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "http://www.com/posts/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	mux.ServeHTTP(w, req)
	assert.Equal(t, "DELETE, GET, POST", w.Header().Get("Allow"))
}

func TestNamedParamValues(t *testing.T) {
	var h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		fmt.Fprintf(w, "post_id=%s&id=%s", q.Get(":post_id"), q.Get(":id"))
	})
	mux := New()
	mux.Add("/posts/:post_id/tags/:id", h, "GET")

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://www.com/posts/123/tags/456", nil)
	if err != nil {
		t.Fatal(err)
	}

	mux.ServeHTTP(w, req)
	assert.Equal(t, "post_id=123&id=456", w.Body.String())
}
