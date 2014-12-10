# mountn

[![Build Status](https://travis-ci.org/nowk/mountn.svg?branch=master)](https://travis-ci.org/nowk/mountn)
[![GoDoc](https://godoc.org/github.com/nowk/mountn?status.svg)](http://godoc.org/github.com/nowk/mountn)

A simple pattern match mux

## Example

    mux := mountn.New()
    mux.Get("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")

      // ...
    }))

    mux.Post("/posts", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // ...
    }))

    mux.Add("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // ...
    }), "PUT", "PATCH")


## License

MIT

