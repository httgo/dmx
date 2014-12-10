package mountn

import (
	"fmt"
	"github.com/nowk/urlp"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type resource struct {
	pat string
	m   urlp.Matcher
	h   http.Handler
}

// NewResource returns a resource with a preconditioned matcher
func NewResource(pat string, h http.Handler) *resource {
	return &resource{
		pat: pat,
		m:   urlp.NewMatcher(pat),
		h:   h,
	}
}

// resources are collection of resources bound to a method
type resources map[string][]*resource

// Add adds a resource bound to a method + pattern. Returns an error if a
// duplicate method + pattern already exists
func (r resources) Add(meth, pat string, h http.Handler) error {
	m, ok := r[meth]
	if ok {
		for _, v := range m {
			if v.pat != pat {
				continue
			}

			return fmt.Errorf("error: mux: %s %s is already defined", meth, v.pat)
		}
	}

	r[meth] = append(r[meth], NewResource(pat, h))

	return nil
}

// mux structure
type mux struct {
	handlers resources
}

func New() *mux {
	return &mux{
		handlers: make(resources),
	}
}

// trim takes off the trailing /
func trim(s string) string {
	s = strings.TrimRight(s, "/")
	if s == "" {
		return "/"
	}

	return s
}

// Add a handler bound to a path pattern and one or more methods. It panics if
// adding a handler returns an error due to a duplication method + pattern
func (m *mux) Add(pat string, h http.Handler, meth ...string) {
	for _, v := range meth {
		err := m.handlers.Add(v, trim(pat), h)
		if err != nil {
			panic(err)
		}
	}
}

func params(p []string, u *url.URL) {
	l := len(p)
	for i, _ := range p {
		n := i + 1
		if n >= l {
			break
		}
		u.RawQuery += "&" + p[i] + "=" + p[n]
	}
}

// match matches a path to a resource's path pattern
func match(r []*resource, u *url.URL) (*resource, bool) {
	for _, v := range r {
		p, ok := v.m.Match(u.Path)
		if !ok {
			continue
		}
		params(p, u)

		return v, ok
	}

	return nil, false
}

func (m *mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r, ok := m.handlers[req.Method]
	if !ok {
		c := 404

		meths, ok := allowed(m.handlers, req)
		if ok {
			c = 405
			w.Header().Add("Allow", strings.Join(meths, ", "))
		}

		http.Error(w, http.StatusText(c), c)
		return
	}

	res, ok := match(r, req.URL)
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	res.h.ServeHTTP(w, req)
}

// allowed returns methods for resources that match the current request path
func allowed(h resources, req *http.Request) ([]string, bool) {
	var meths []string
	for k, v := range h {
		if k == req.Method {
			continue
		}

		_, ok := match(v, req.URL)
		if ok {
			meths = append(meths, k)
		}
	}
	if len(meths) == 0 {
		return nil, false
	}

	sort.Strings(meths)

	return meths, true
}

func (m *mux) Get(pat string, h http.Handler) {
	m.Add(pat, h, "GET")
}

// Geth registers both a head and get handler
func (m *mux) Geth(pat string, h http.Handler) {
	m.Add(pat, h, "HEAD", "GET")
}

func (m *mux) Head(pat string, h http.Handler) {
	m.Add(pat, h, "HEAD")
}

func (m *mux) Post(pat string, h http.Handler) {
	m.Add(pat, h, "POST")
}

func (m *mux) Put(pat string, h http.Handler) {
	m.Add(pat, h, "PUT")
}

// Putp registers both a put and patch handler
func (m *mux) Putp(pat string, h http.Handler) {
	m.Add(pat, h, "PUT", "PATCH")
}

func (m *mux) Patch(pat string, h http.Handler) {
	m.Add(pat, h, "PATCH")
}

func (m *mux) Del(pat string, h http.Handler) {
	m.Add(pat, h, "DELETE")
}
