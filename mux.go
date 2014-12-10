package dmx

import (
	"fmt"
	"github.com/nowk/urlp"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type resource struct {
	http.Handler
	urlp.Matcher

	pat string
}

// NewResource returns a resource with a preconditioned matcher from the pattern
func NewResource(pat string, h http.Handler) *resource {
	return &resource{
		h,
		urlp.NewMatcher(pat),
		pat,
	}
}

// mux is a collection of method bound resources
type mux map[string][]*resource

func New() mux {
	return make(mux)
}

// add adds a new resource given a single method, patter and handler. Returning
// and error on a pattern + method duplication
func (r mux) add(meth, pat string, h http.Handler) error {
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
func (m mux) Add(pat string, h http.Handler, meth ...string) {
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
		p, ok := v.Match(pathStr)
		if ok {
			return v, p, ok
		}
	}

	return nil, nil, false
}

// ServeHTTP implements http.Handler
func (m mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r, ok := m[req.Method]
	if !ok {
		c := 404
		if meths, ok := methodsAllowed(m, req); ok {
			c = 405
			w.Header().Add("Allow", strings.Join(meths, ", "))
		}
		http.Error(w, http.StatusText(c), c)
		return
	}

	res, p, ok := Match(r, req.URL.Path)
	if !ok {
		http.Error(w, http.StatusText(404), 404)
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

func methodsAllowed(m mux, req *http.Request) ([]string, bool) {
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
