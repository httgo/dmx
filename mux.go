package dmx

import (
	"fmt"
	"gopkg.in/nowk/urlp.v1"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type resource struct {
	http.Handler

	pat string
}

// NewResource returns a resource with a preconditioned matcher from the pattern
func NewResource(pat string, h http.Handler) *resource {
	return &resource{
		h,
		pat,
	}
}

// Mux is a collection of method bound resources
type Mux map[string][]*resource

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

// trim trims the trailing slash. Will always return atleast "/"
func trim(s string) string {
	s = strings.TrimRight(s, "/")
	if s == "" {
		return "/"
	}

	return s
}

func Match(r []*resource, pathStr string) (*resource, []string, bool) {
	for _, v := range r {
		p, ok := urlp.Match(v.pat, pathStr)
		if ok {
			return v, p, ok
		}
	}

	return nil, nil, false
}

// notFound handles 404 and 405 errors looking up the path in other method sets
// and returns an Allow header if the path is allowed on other methods
func (m Mux) notFound(w http.ResponseWriter, req *http.Request) {
	c := 404

	meths, ok := methodsAllowed(m, req)
	if ok {
		c = 405
		w.Header().Add("Allow", strings.Join(meths, ", "))
	}

	http.Error(w, http.StatusText(c), c)
}

// ServeHTTP implements http.Handler
	r, ok := m[req.Method]
	if !ok {
		m.notFound(w, req)
		return
	}

	res, p, ok := Match(r, req.URL.Path)
func (m Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !ok {
		m.notFound(w, req)
		return
	}

	params(p, req.URL)
	res.ServeHTTP(w, req)
}

func params(p []string, u *url.URL) {
	i := 0
	l := len(p)
	for i < l {
		u.RawQuery = u.RawQuery + "&" + p[i] + "=" + p[i+1]
		i = i + 2
	}
}

func methodsAllowed(m Mux, req *http.Request) ([]string, bool) {
	var meths []string
	for k, v := range m {
		if k != req.Method {
			_, _, ok := Match(v, req.URL.Path)
			if ok {
				meths = append(meths, k)
			}
		}
	}

	if len(meths) == 0 {
		return nil, false
	}

	sort.Strings(meths)
	return meths, true
}

func (m *Mux) Get(pat string, h http.Handler) {
	m.Add(pat, h, "GET")
}

// Geth registers both a head and get handler
func (m *Mux) Geth(pat string, h http.Handler) {
	m.Add(pat, h, "HEAD", "GET")
}

func (m *Mux) Head(pat string, h http.Handler) {
	m.Add(pat, h, "HEAD")
}

func (m *Mux) Post(pat string, h http.Handler) {
	m.Add(pat, h, "POST")
}

func (m *Mux) Put(pat string, h http.Handler) {
	m.Add(pat, h, "PUT")
}

// Putp registers both a put and patch handler
func (m *Mux) Putp(pat string, h http.Handler) {
	m.Add(pat, h, "PUT", "PATCH")
}

func (m *Mux) Patch(pat string, h http.Handler) {
	m.Add(pat, h, "PATCH")
}

func (m *Mux) Del(pat string, h http.Handler) {
	m.Add(pat, h, "DELETE")
}
