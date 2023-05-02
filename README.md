# framework-muxer-showdown

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## net/http implementation

### General considerations

It's pretty powerful out of the box. Dealing with `http.Handler` and `http.HandlerFunc` is actually easier than I anticipated. Creating middlewares is easy and fairly straightforward with them being functions that take a handler interface, and return a handler interface.

Nested muxers are also very very easy to do, though we do need to keep a few things in mind when nesting them.

Error handling is easy, once modifying the request context is done, and middlewares are implemented.

Path variables are a pain.

Here are the details of each stuff:

### Details of criteria

#### Context type

This is an embedded standar library `context.Context`. We can get it with `request.Context()`. See the error handling section on how to modify the request in a way that the changed context propagates to other parts of the codebase.

#### Standard library handling

It is standard library, out of the box. :)

#### Accessing raw Request and ResponseWriter

`net/http` gives us really two ways of dealing with handlers. Either it's an `http.HandlerFunc`, or whatever implements the `http.Handler` interface. The latter has a `ServeHTTP` method, which itself is an `http.HandlerFunc`, so the two of them are mostly interchangeable. The signature for the func is
```go
func thing(w http.ResponseWriter, r *http.Request) {}
```
So yep, this is about as raw request as it gets.

Routing is done via a muxer (multiplexer) like this:
```go
mux := http.NewServeMux()

mux.Handle("/somepath", handlers.SomeHandlerFunc)
mux.Handle("/otherpath", logic.ImplementsHandlerInterface)
```

#### Websocket

I tried two ways of implementing websockets, one using the standard library websocket implementation at
`"golang.org/x/net/websocket"`, the other using gobwas from `"github.com/gobwas/ws"` because initially I had a problem with getting the standard library websocket working as I kept getting 401.

##### Using standard library

The handler needs to be `websocket.Handler`, that one takes a connection. That's on the handler side, so inside all we really need is a for loop with receive and send bits.

On the muxer side, we don't need to do anything, because the `websocket.Handler` implements the `http.Handler` interface, so we can just plug it in.

The incoming request needs to have an `Origin` header present with a parseable URL as value, so in real world uses that won't be a problem, but Postman might catch you out in this.

See `handlers.WSStd()` for how that works.

##### Using gobwas

We need to upgrade the request to a connecting using `ws.Upgrade`, and then it's a standard for loop. See `handlers.WS()` for how it works.

#### Path specificity

This is a super yes.

* `/single` matches a single path: `/single`. And done, that's it.
* `/single/` matches all paths that have `/single/` as a prefix, so `/single/`, `/single/bla`, `/single/ladies/put/a/ring/on/it`
* `/single/long/thing` however matches a specific handler, even though it should also be covered by the `/single/` matcher. This is more specific, so matches first.

#### Path variables

This is... eeeehhh... it _can_ be done, but it needs a LOT of work. There's this [super useful article by Ben Hoyt about different ways of handling dynamic routing](https://benhoyt.com/writings/go-routing/), which is easy to understand, and can work, but we'd need to actually build it ourselves, and this is why routers exist.

#### Grouping

This was super easy. Take a muxer, and put another muxer onto a path. The only important thing that we need to pay attention to are slashes and using `http.StripPrefix` like so:

```go
// Grouping
groupMux := http.NewServeMux()
// will handle /v1/hello
groupMux.Handle("/hello", getMiddleware(handlers.Hello()))

mux.Handle("/v1/", http.StripPrefix("/v1", groupMux))
```
The group prefix NEEDS to have the trailing slash on it, the `http.StripPrefix` needs to NOT have the trailing slash, and the paths in the nested muxer need to have the trailing slash. Otherwise it just works.

#### Overlaps

This is the situation where if we have a handler for a dynamic routing, like `/page/+pageslug` and `/page/contact`, because technically `contact` could also be a dynamic thing, but because path variables with `net/http` is HARD, I didn't check.

#### General middleware

Fairly easy, Middlewares look like this: `func(next http.Handler) http.Handler`, and then within them we call `next.ServeHTTP(w, r)`, and ... done.

The one interesting thing I learned is that returning a `http.Handler` as a function you can do this:
```go
func SomeMiddleware(l zerolog.Logger, db package.Struct) func(http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do the thing of the middleware here
        l.Info().Msg("doing middleware work")

        next.ServeHTTP(w, r)
    })
}
```
And then in the muxer we can do the following
```go
mux := http.NewServeMux()

// option 1
mux.Handle("/some-path", SomeMiddleware(l, db)(handlers.WhateverHandler))

// option 2
mw := SomeMiddleware(l, db)
mux.Handle("/some-path", mw(handlers.WhateverHandler))
```

#### Error handling middleware

Eh, so this is kinda. Technically we can shove errors into the request ctx, which I added a `web` package helper functions to do it for us, and then in an error handler middleware we can check the ctx, and act accordingly.

The reason it's a `Kinda` and not a `Yes` is because in order to change the context on a pointer to a request, you need to do the following at the end of the handler:
```go
*r = *r.Clone(newCtx)
```
Yes, we assign a pointer to a pointer. And then dealing with error type checking using pointers is also kind of mindboggling, but otherwise it works as expected. It's just clunky.

#### Context up and down

As with the error handling, it can be done, but pointers = pointers, `request.Clone`, `request.WithContext`, but ... can be done. It's `Kinda`, because it's clunky.

By up and down, I mean passing context from a middleware into a nested next handler (down), and grabbing a changed context from within a nested handler back to a middleware (up).

#### Unit tests

Super easy to do! Standar library's `httptest` package seamlessly integrates with the handlers. The tests are easy to set up and easy to understand, and do not actually need the app to be started, whether in docker container, or in a different process, so that's excellent!

Moreover we can absolutely mock any dependant services, like database, loggers, tracers, etc, which would make life significantly easier too!

See the code at [handlers/errors_test.go](handlers/errors_test.go)!

#### Ecosystem

There's not much to talk about here. Pretty much everything is handmade. Though the examples for most things, like loggers etc, are given in standard library, so copying-pasting should be readily available.

#### Performance
Standard library muxer does not have any router configurable panic or error handlers.

##### /performance
* [perftests/run1.log](perftests/run1.log):  Requests/sec:	25386.6980
* [perftests/run2.log](perftests/run2.log):  Requests/sec:	24801.1679
* [perftests/run3.log](perftests/run3.log):  Requests/sec:	24726.3641
* [perftests/run4.log](perftests/run4.log):  Requests/sec:	24558.7190

##### /smol-perf
* [perftests/smol-run1.log](perftests/smol-run1.log):  Requests/sec:	65754.6579
* [perftests/smol-run2.log](perftests/smol-run2.log):  Requests/sec:	54035.5844
* [perftests/smol-run3.log](perftests/smol-run3.log):  Requests/sec:	56846.3588
* [perftests/smol-run4.log](perftests/smol-run4.log):  Requests/sec:	56624.9325
