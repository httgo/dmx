package dmx

import (
	"net/http"
)

// Extends the Mux struct with basic short methods

func (m *Mux) Get(path string, w ...interface{}) *Mux {
	return m.Add("GET", path, w...)
}

func (m *Mux) Post(path string, w ...interface{}) *Mux {
	return m.Add("POST", path, w...)
}

func (m *Mux) Put(path string, w ...interface{}) *Mux {
	return m.Add("PUT", path, w...)
}

func (m *Mux) Patch(path string, w ...interface{}) *Mux {
	return m.Add("PATCH", path, w...)
}

func (m *Mux) Delete(path string, w ...interface{}) *Mux {
	return m.Add("DELETE", path, w...)
}

func (m *Mux) Head(path string, w ...interface{}) *Mux {
	return m.Add("HEAD", path, w...)
}

func (m *Mux) Options(path string, w ...interface{}) *Mux {
	return m.Add("OPTIONS", path, w...)
}

func (m *Mux) GetFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Get(path, http.HandlerFunc(fn))
}

func (m *Mux) PostFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Post(path, http.HandlerFunc(fn))
}

func (m *Mux) PutFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Put(path, http.HandlerFunc(fn))
}

func (m *Mux) PatchFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Patch(path, http.HandlerFunc(fn))
}

func (m *Mux) DeleteFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Delete(path, http.HandlerFunc(fn))
}

func (m *Mux) HeadFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Head(path, http.HandlerFunc(fn))
}

func (m *Mux) OptionsFunc(
	path string, fn func(http.ResponseWriter, *http.Request)) *Mux {

	return m.Options(path, http.HandlerFunc(fn))
}
