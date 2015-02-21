package dmx

import (
	"net/http"
)

// Mux is a collection of method bound resources
type Mux map[string]resources

func New() Mux {
	return make(Mux)
}

// Add adds a new resource given the pattern, handler and one or more methods.
// Panics on a pattern + method duplication
func (m Mux) Add(pat string, h http.Handler, meths ...string) {
	res := NewResource(meths, pat, h)
	err := res.Apply(m)
	if err != nil {
		panic(err)
	}
}

// Then returns the final serve handler through an Alice style constructor
// allowing you to passing in your own NotFound handler. Or push use `Then` to
// stack the mux into a middleware chain
func (m Mux) Then(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if res, ok := Match(m, req); ok {
			res.ServeHTTP(w, req)
			return
		}

		h.ServeHTTP(w, req)
	})
}

// Match returns a matching resources based on a matching pattern to path and
// request method
func Match(m Mux, req *http.Request) (*resource, bool) {
	if r, ok := m[req.Method]; ok {
		return r.Match(req)
	}

	return nil, false
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
