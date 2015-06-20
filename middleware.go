package dmx

import (
	"net/http"
)

// Middleware represents the interface for a middleware
type Middleware interface {
	Then(http.Handler) http.Handler
}

// MiddlewareFunc implements Middleware
type MiddlewareFunc func(http.Handler) http.Handler

func (m MiddlewareFunc) Then(next http.Handler) http.Handler {
	return m(next)
}
