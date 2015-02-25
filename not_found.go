package dmx

import (
	"net/http"
	"strings"
)

// NotFound types a mux to handle 404s and 405s returning an Allow header if
// available
type NotFound Mux

func (n NotFound) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := 404
	meths, ok := n.AllowedMethods(req)
	if ok {
		c = 405
		w.Header().Add("Allow", strings.Join(meths, ", "))
	}

	http.Error(w, http.StatusText(c), c)
}

func (n NotFound) AllowedMethods(req *http.Request) ([]string, bool) {
	var meths []string
	for k, v := range n {
		if k == req.Method {
			continue
		}

		if r := v.Match(req); r != nil {
			meths = append(meths, k)
		}
	}
	if len(meths) == 0 {
		return nil, false
	}
	return meths, true
}
