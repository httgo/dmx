package dmx

import (
	"errors"
	"net/http"
)

// Mux is the router structure for middlewares and resources
type Mux struct {
	http.Handler

	// resources is a collection of method bound handlers
	resources resources

	musts       []Middleware
	middlewares []Middleware
	mounts      []Middleware

	// next is the handler called when no resource for a request can be matched,
	// this by default will use http.NotFound
	next http.Handler
}

func New(opts ...func(*Mux)) *Mux {
	m := &Mux{
		resources: make(resources),
	}

	m.Else(http.HandlerFunc(http.NotFound))

	for _, v := range opts {
		v(m)
	}

	return m
}

// Else sets next and reconstructs Handler via Then(next)
func (m *Mux) Else(next http.Handler) http.Handler {
	m.next = next

	return m.Then(m.next)
}

func construct(m []Middleware, h http.Handler) http.Handler {
	i := len(m) - 1
	for ; i >= 0; i-- {
		h = m[i].Then(h)
	}

	return h
}

// Then implements middleware as well as defining the internal Handler field
func (m *Mux) Then(next http.Handler) http.Handler {
	m.Handler = m.resources.Then(next)

	if len(m.musts) > 0 {
		m.Handler = construct(m.musts, m.Handler)
	}

	return m.Handler
}

func checkMiddleware(v interface{}) Middleware {
	h, ok := v.(func(http.Handler) http.Handler)
	if ok {
		return MiddlewareFunc(h)
	}

	m, ok := v.(Middleware)
	if ok {
		return m
	}

	panic(errors.New("invalid middleware"))

	return nil
}

// Must allows you to append middlewares that MUST be run before any matching
// resources are looked up (regardless of whether the resource was found or
// not). Must middlewares are only executed upon entry of the Mux. Mounts with
// Must midddlewares will only execute those middlewares if the Mount is entered
// to search for a resource.
func (m *Mux) Must(v interface{}) *Mux {
	m.musts = append(m.musts, checkMiddleware(v))

	m.Then(m.next) // must rebuild Handler

	return m
}

// Use appends middlewares for the current Mux. Middlewares are only executed
// for matched resources. If you require a Middleware to always be executed
// use Must(...).
//
// This must be called before applying any resources or Mounts
func (m *Mux) Use(v interface{}) *Mux {
	m.middlewares = append(m.middlewares, checkMiddleware(v))

	return m
}

func extend(m []Middleware, mnt *Mux) *Mux {
	for _, u := range mnt.resources {
		for _, v := range u {
			v.Prefix(m...)
		}
	}

	for _, v := range mnt.mounts {
		extend(m, v.(*Mux)) // recurse down mount tree
	}

	return mnt
}

func (m *Mux) Mount(mnt *Mux) *Mux {
	m.mounts = append(m.mounts, extend(m.middlewares, mnt))

	h := construct(m.mounts[1:], m.mounts[0].(http.Handler))
	m.Then(h)

	return m
}

func (m *Mux) Add(meth, path string, v ...interface{}) *Mux {
	w := make([]interface{}, 0, len(m.middlewares))
	for _, v := range m.middlewares {
		w = append(w, v)
	}

	m.resources[meth] = append(
		m.resources[meth], newResource(meth, path, append(w, v...)...))

	return m
}
