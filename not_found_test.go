package dmx

import (
	"gopkg.in/nowk/assert.v2"
	"net/http"
	"sort"
	"testing"
)

func TestAllowedMethods(t *testing.T) {
	mux := New()
	mux.Add("/posts", hfunc(""), "GET")
	mux.Add("/posts/:id", hfunc(""), "PUT", "POST", "DELETE")

	req, err := http.NewRequest("GET", "http://www.com/posts/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := NotFound(mux).AllowedMethods(req)
	sort.Strings(m)
	assert.True(t, ok)
	assert.Equal(t, []string{"DELETE", "POST", "PUT"}, m)
}
