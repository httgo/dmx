package dmx

import (
	"fmt"
	"net/http"
	"strings"
)

type resource struct {
	http.Handler

	// Methods are the HTTP methods supported by this resource
	Methods []string

	// Pattern is the path pattern for this resource
	Pattern string
}

func NewResource(meths []string, pat string, h http.Handler) *resource {
	return &resource{
		Handler: h,
		Methods: meths,
		Pattern: trim(pat),
	}
}

func (r *resource) Apply(mux *Mux) error {
	for _, m := range r.Methods {
		p, ok := mux.Resources[m]
		if ok {
			for _, v := range p {
				if v.Pattern == r.Pattern {
					return fmt.Errorf("error: mux: %s %s is already defined", m, r.Pattern)
				}
			}
		}

		mux.Resources[m] = append(mux.Resources[m], r)
	}
	return nil
}

// trim trims the trailing slash. Will always return atleast "/"
func trim(s string) string {
	s = strings.TrimRight(s, "/")
	if s == "" {
		return "/"
	}
	return s
}
