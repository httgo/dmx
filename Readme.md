# dmx

[![Build Status](https://travis-ci.org/httgo/dmx.svg?branch=master)](https://travis-ci.org/httgo/dmx)
[![GoDoc](https://godoc.org/gopkg.in/httgo/dmx.v4?status.svg)](http://godoc.org/gopkg.in/httgo/dmx.v4)

A simple pattern match mux.


## Install

    go get gopkg.in/httgo/dmx.v4


## Usage

    var getPostHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")

      // ... 
    })

    func main() {
      mux := dmx.New()
      mux.Get("/posts/:id", getPostHandler)

      err := http.ListenAndServe(":3000", mux)
      if err != nil {
        log.Fatalf("listen: %s", err)
      }
    }

---

__Use(Middleware)__

    mux.Use(func(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        next.ServerHTTP(w, req)
      })
    })

Use only get called for a matching resource. If you need a middleware to run on all requests use `Must`.

---

__Must(Middleware)__

    mux.Must(func(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        next.ServerHTTP(w, req)
      })
    })

Must executes on all requests.


## License

MIT

