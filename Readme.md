# dmx

[![Build Status](https://travis-ci.org/httgo/dmx.svg?branch=master)](https://travis-ci.org/httgo/dmx)
[![GoDoc](https://godoc.org/gopkg.in/httgo/dmx.v2?status.svg)](http://godoc.org/gopkg.in/httgo/dmx.v2)

A simple pattern match mux. *A speed experiment.*


## Install

    go get gopkg.in/httgo/dmx.v2


## Example

    package main

    import "net/http"
    import "gopkg.in/httgo/dmx.v2"

    var getPostHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")

      // ... 
    })

    func main() {
      mux := dmx.New()
      mux.Get("/posts/:id", getPostHandler)

      err := http.ListenAndServe(":3000", mux.Then(dmx.NotFound(mux)))
      if err != nil {
        log.Fatalf("fatal: listen: %s", err)
      }
    }

##### Basic Method Shortcuts

GET

    mux.Get(string, http.Handler)
    mux.GetFunc(string, http.HandlerFunc)
    
POST
    
    mux.Post(string, http.Handler)
    mux.PostFunc(string, http.HandlerFunc)
    
PUT

    mux.Put(string, http.Handler)
    mux.PutFunc(string, http.HandlerFunc)
    
DELETE

    mux.Del(string, http.Handler)
    mux.DelFunc(string, http.HandlerFunc)

---

Use the `Add` method for binding multiple methods to a path and handler.

    mux.Add("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // ...
    }), "PUT", "PATCH")

---

`.:format` at the end of your path pattern will parse the extension and provide it as a parameter named `:_format`.

    mux.Get("/posts/:id.:format", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")
      format := v.Get(":_format")

      // ...
    }))

Using `.:format` will match paths with or without the extension.

---

##### Handling Not Found

You define how you want non-matches to be handled. 

You can use the built in `dmx.NotFound`. This will return either `404` or `405` (with an `Allow` header).

    h := mux.Then(dmx.NotFound(mux))

Or you you can pass it off to another handler. Ex: a file server

    stc := http.Dir("./public")
    pub := http.FileServer(stc)
    ...
    h := mux.Then(pub)

*Note: You must define a `Then` handler, the mux itself does not implement `http.Handler`*


## License

MIT

