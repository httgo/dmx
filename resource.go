package dmx

import (
	"net/http"
	"net/url"

	"gopkg.in/nowk/urlp.v2"
)

type resources map[string][]*resource

func (r resources) Then(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h, ok := r.Match(req)
		if ok {
			h.ServeHTTP(w, req)

			return
		}

		next.ServeHTTP(w, req)
	})
}

// params attaches the url parameters to url's query string
func params(p []string, u *url.URL) {
	n := len(p)
	for i := 0; i < n; {
		u.RawQuery = u.RawQuery + "&" + p[i] + "=" + p[i+1]
		i = i + 2
	}
}

func (r resources) Match(req *http.Request) (http.Handler, bool) {
	m, ok := r[req.Method]
	if !ok {
		return nil, false
	}

	for _, v := range m {
		p, ok := urlp.Match(v.Path, req.URL.Path)
		if !ok {
			continue
		}

		params(p, req.URL)

		return v, true
	}

	return nil, false
}

type resource struct {
	http.Handler

	Method string
	Path   *urlp.Path
}

func parsemh(v ...interface{}) ([]Middleware, http.Handler) {
	n := len(v) - 1
	m := make([]Middleware, 0, n+1)
	for _, w := range v[:n] {
		m = append(m, w.(Middleware))
	}

	return m, v[n].(http.Handler)
}

func newResource(meth, path string, v ...interface{}) *resource {
	m, h := parsemh(v...)

	return &resource{
		Handler: construct(m, h),
		Method:  meth,
		Path:    urlp.NewPath(path),
	}
}

func (r *resource) Prefix(m ...Middleware) {
	r.Handler = construct(m, r.Handler)
}
