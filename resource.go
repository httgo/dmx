package dmx

import (
	"gopkg.in/nowk/urlp.v1"
	"net/http"
	"net/url"
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

func params(p []string, u *url.URL) {
	i := 0
	l := len(p)
	for i < l {
		u.RawQuery = u.RawQuery + "&" + p[i] + "=" + p[i+1]
		i = i + 2
	}
}
