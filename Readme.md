# dmx

[![Build Status](https://travis-ci.org/nowk/dmx.svg?branch=master)](https://travis-ci.org/nowk/dmx)
[![GoDoc](https://godoc.org/github.com/nowk/dmx?status.svg)](http://godoc.org/github.com/nowk/dmx)

A simple pattern match mux. *A speed experiment.*

## Example

    mux := dmx.New()
    mux.Get("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")

      // ...
    }))

    mux.Post("/posts", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // ...
    }))

Handling multple methods

    mux.Add("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // ...
    }), "PUT", "PATCH")


## License

MIT

