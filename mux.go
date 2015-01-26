package dmx

import (
	"fmt"
	"gopkg.in/nowk/urlp.v1"
	"net/http"
	"net/url"
	"strings"
)

type resource struct {
	http.Handler
	pat string
}

func NewResource(pat string, h http.Handler) *resource {
	return &resource{
		Handler: h,
		pat:     pat,
	}
}

type resources []*resource

func params(p []string, u *url.URL) {
	i := 0
	l := len(p)
	for i < l {
		u.RawQuery = u.RawQuery + "&" + p[i] + "=" + p[i+1]
		i = i + 2
	}
}

// Match returns a matching resources based on a matching pattern to path
func (r resources) Match(req *http.Request) (*resource, bool) {
	for _, v := range r {
		p, ok := urlp.Match(v.pat, req.URL.Path)
		if ok {
			params(p, req.URL)
			return v, ok
		}
	}
	return nil, false
}

// Mux is a collection of method bound resources
type Mux map[string]resources

func New() Mux {
	return make(Mux)
}

// add adds a new resource given a single method, patter and handler. Returning
// and error on a pattern + method duplication
func (r Mux) add(meth, pat string, h http.Handler) error {
	m, ok := r[meth]
	if ok {
		for _, v := range m {
			if v.pat == pat {
				return fmt.Errorf("error: mux: %s %s is already defined", meth, v.pat)
			}
		}
	}

	r[meth] = append(r[meth], NewResource(pat, h))
	return nil
}

// trim trims the trailing slash. Will always return atleast "/"
func trim(s string) string {
	s = strings.TrimRight(s, "/")
	if s == "" {
		return "/"
	}

	return s
}

// Add adds a new resource given the pattern, handler and one or more methods.
// Panics on a pattern + method duplication
func (m Mux) Add(pat string, h http.Handler, meth ...string) {
	for _, v := range meth {
		err := m.add(v, trim(pat), h)
		if err != nil {
			panic(err)
		}
	}
}

// Handler returns the final handler delegating any unmatched routes to the
// provided handler.
func (m Mux) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res, ok := Match(m, req)
		if !ok {
			h.ServeHTTP(w, req)
			return
		}

		res.ServeHTTP(w, req)
	})
}

// Match returns a matching resources based on a matching pattern to path and
// request method
func Match(m Mux, req *http.Request) (*resource, bool) {
	r, ok := m[req.Method]
	if !ok {
		return nil, false
	}

	return r.Match(req)
}

func (m *Mux) Get(pat string, h http.Handler) {
	m.Add(pat, h, "GET")
}

func (m *Mux) GetFunc(pat string, fn http.HandlerFunc) {
	m.Get(pat, http.HandlerFunc(fn))
}

// Geth registers both a head and get handler
func (m *Mux) Geth(pat string, h http.Handler) {
	m.Add(pat, h, "HEAD", "GET")
}

func (m *Mux) Head(pat string, h http.Handler) {
	m.Add(pat, h, "HEAD")
}

func (m *Mux) HeadFunc(pat string, fn http.HandlerFunc) {
	m.Head(pat, http.HandlerFunc(fn))
}

func (m *Mux) Post(pat string, h http.Handler) {
	m.Add(pat, h, "POST")
}

func (m *Mux) PostFunc(pat string, fn http.HandlerFunc) {
	m.Post(pat, http.HandlerFunc(fn))
}

func (m *Mux) Put(pat string, h http.Handler) {
	m.Add(pat, h, "PUT")
}

func (m *Mux) PutFunc(pat string, fn http.HandlerFunc) {
	m.Put(pat, http.HandlerFunc(fn))
}

// Putp registers both a put and patch handler
func (m *Mux) Putp(pat string, h http.Handler) {
	m.Add(pat, h, "PUT", "PATCH")
}

func (m *Mux) Patch(pat string, h http.Handler) {
	m.Add(pat, h, "PATCH")
}

func (m *Mux) PatchFunc(pat string, fn http.HandlerFunc) {
	m.Patch(pat, http.HandlerFunc(fn))
}

func (m *Mux) Del(pat string, h http.Handler) {
	m.Add(pat, h, "DELETE")
}

func (m *Mux) DelFunc(pat string, fn http.HandlerFunc) {
	m.Del(pat, http.HandlerFunc(fn))
}
