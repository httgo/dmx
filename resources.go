package dmx

import (
	"gopkg.in/nowk/urlp.v1"
	"net/http"
	"net/url"
)

type resources []*resource

// Match returns a matching resources based on a matching pattern to path
func (r resources) Match(req *http.Request) (*resource, bool) {
	for _, v := range r {
		p, ok := urlp.Match(v.Pattern, req.URL.Path)
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
