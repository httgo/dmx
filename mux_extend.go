package dmx

import (
	"net/http"
)

// Extends the Mux struct with basic short methods

func (m *Mux) GET(path string, w ...interface{}) *Mux {
	return m.Add("GET", path, w...)
}

func (m *Mux) POST(path string, w ...interface{}) *Mux {
	return m.Add("POST", path, w...)
}

func (m *Mux) PUT(path string, w ...interface{}) *Mux {
	return m.Add("PUT", path, w...)
}

func (m *Mux) PATCH(path string, w ...interface{}) *Mux {
	return m.Add("PATCH", path, w...)
}

func (m *Mux) DELETE(path string, w ...interface{}) *Mux {
	return m.Add("DELETE", path, w...)
}

func (m *Mux) HEAD(path string, w ...interface{}) *Mux {
	return m.Add("HEAD", path, w...)
}

func (m *Mux) GETFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.GET(path, http.HandlerFunc(fn))
}

func (m *Mux) POSTFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.POST(path, http.HandlerFunc(fn))
}

func (m *Mux) PUTFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.PUT(path, http.HandlerFunc(fn))
}

func (m *Mux) PATCHFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.PATCH(path, http.HandlerFunc(fn))
}

func (m *Mux) DELETEFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.DELETE(path, http.HandlerFunc(fn))
}

func (m *Mux) HEADFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.HEAD(path, http.HandlerFunc(fn))
}
