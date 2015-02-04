package dmx

import (
	"net/http"
	"sort"
	"strings"
)

// NotFound types a mux to handle 404s and 405s returning an Allow header if
// available
type NotFound struct {
	*Mux
}

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
	for k, v := range n.Resources {
		if k != req.Method {
			_, ok := v.Match(req)
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
