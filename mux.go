package dmx

import (
	"net/http"
)

// Mux is a collection of method bound resources
type Mux struct {
	Resources       map[string]resources
	NotFoundHandler http.Handler
}

func New() *Mux {
	var m Mux
	m = Mux{
		Resources:       make(map[string]resources),
		NotFoundHandler: NotFound{&m},
	}
	return &m
}

// Add adds a new resource given the pattern, handler and one or more methods.
// Panics on a pattern + method duplication
func (m *Mux) Add(pat string, h http.Handler, meths ...string) {
	r := NewResource(meths, pat, h)
	if err := r.Apply(m); err != nil {
		panic(err)
	}
}

func (m Mux) NotFound(w http.ResponseWriter, req *http.Request) {
	m.NotFoundHandler.ServeHTTP(w, req)
}

func (m Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res, ok := Match(&m, req)
	if !ok {
		m.NotFound(w, req)
		return
	}
	res.ServeHTTP(w, req)
}

// Match returns a matching resources based on a matching pattern to path and
// request method
func Match(m *Mux, req *http.Request) (*resource, bool) {
	r, ok := m.Resources[req.Method]
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
