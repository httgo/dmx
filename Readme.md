# dmx

[![Build Status](https://travis-ci.org/httgo/dmx.svg?branch=master)](https://travis-ci.org/httgo/dmx)
[![GoDoc](https://godoc.org/gopkg.in/httgo/dmx.v3?status.svg)](http://godoc.org/gopkg.in/httgo/dmx.v3)

A simple pattern match mux. *A speed experiment.*


## Install

    go get gopkg.in/httgo/dmx.v3


## Usage

    package main

    import "net/http"
    import "gopkg.in/httgo/dmx.v3"

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

*Note: You must call `Then`, the mux itself does not implement `http.Handler`*

---

#### Mounting other Muxes

You can mount other Muxes onto a parent mux, providing a way to organize large mux structures.

    posts := dmx.New()
    posts.Get("/posts", ...)
    posts.Get("/posts/:id", ...)

    mux := dmx.New()
    mux.Get("/", ...)
    mux.Mount(posts)

This would be equivalent to

    mux := dmx.New()
    mux.Get("/", ...)
    mux.Get("/posts", ...)
    mux.Get("/posts/:id", ...)

You can also mount on a namespace

    posts := dmx.New()
    posts.Get("/", ...)
    posts.Get("/:id", ...)

    mux := dmx.New()
    mux.MountAt("/posts", posts)


## License

MIT

