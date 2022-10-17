# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## chi implementation

### General Considerations

chi is weird. The handler functions are `http.HandlerFunc` types, so standard library, so I'm trying to figure out (after only implementing the health endpoint) why we need chi on top of the net/http router at all. Hopefully this is going to be apparent soon.

#### Graceful shutdown

chi is just a muxer rather than an entire framework, which means realistically instead of launching chi itself with
```go
r := chi.NewRouter()

http.ListenAndServe(":3000", r)
```
We should attach it to the standard library server:
```go
r := chi.NewRouter()

server := &http.Server{
    Addr:    ":9000",
    Handler: r,
}

server.ListenAndServer()
```

### Details of criteria

#### Context type

It's the embedded standard context into the *http.Request

#### Standard library handling

That's what the handlers are, so uh... yeah, it supports it.

#### Accessing raw Request and ResponseWriter

Because handlers are just standard library ones, the request and response are freely available, and no need to grab them out of something embedded.

#### Websocket

Very similar situation that the net/http. The only minor niggle is that the `r.Get` function and siblings expect an `http.HandlerFunc`, but the golang `websocket` handler is an `http.Handler` interface type, so instead of using the actual websocket handler, we need to use the `ServeHTTP` method instead (without brackets). Otherwise works as expected.

#### Path specificity

Works the same as net/http and echo. It supports a wildcard (`/spec/*`), but unlike gin, it correctly handles a longer more specific route even if that would conflict with the wildcard.

#### Path variables

Yep, supports them. The declaration is between curly braces, rather than leading colon or plus sign, but otherwise easy to understand, and the [documentation is clear](https://go-chi.io/#/pages/routing?id=routing-patterns-amp-url-parameters).

chi has a bunch of convenience functions that extract the data from the standard library request, so thumbs up!

Additional really cool feature is that we can use a regex pattern to restrict the url params, so for example a route that looks like this:
```
/article/{date}-{slug}
```
and a request to this url:
```
http://server.com/article/2022-10-11-100-awesome-things
```
can be ambiguous, because which `-` is the one that separates the date / slug? Could be either
```
2022 / 10-11-100-awesome-things
2022-10 / 11-100-awesome-things
2022-10-11 / 100-awesome-things
2022-10-11-100 / awesome-things
2022-10-11-100-awesome / things
```
But if the route is declared like this:
```
/articles/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}-{article}
```
The only one that matches that is the following:
```
2022-10-11 / 100-awesome-things
```

#### Grouping
Grouping is super powerful with chi. There are two ways of doing it, using sub routers, or using groups. Both of them can be then mounted on a prefix, as they implement the `http.Handler` interface.

The main difference between sub-routers and groups seems to be that a sub router **needs** to be mounted on a parent router, whereas you can create a bunch of different groups on the root router to segregate different routes to use different middlewares if you don't want to use a custom prefix.

So given the following example routes:
```shell
/contact
/login
/dashboard
/account
/logout
/shop
```
You only want to place an auth middleware on `account` and `dashboard`, but not the others. The two existing solutions so far have been either to attache the middleware to the individual handlers, which is a lot of code duplication, or to change those routes to `/admin/*`, and attach the middleware to the admin group.

With chi's groups you can do this:
```go
r := chi.NewRouter()
r.Get("/contact", handlers.Contact())
r.Get("/login", handlers.Login())
r.Get("/logout", handlers.Logout())
r.Get("/shop", handlers.Shop())

r.Group(func(gr chi.Router) {
	gr.Use(middlewares.Auth())
	gr.Get("/dashboard", handlers.Dashboard())
	gr.Get("/account", handlers.Account())
})
```
This is super neat!

#### Overlaps

Has no problem supporting the use case.

#### General middleware

A middleware is a `func(http.Handler) http.Handler`, same as net/http's case. See the implementations in [middlewares.go](handlers/middlewares.go).

With chi we can either use the `r.Use()` method, or manually wrap the end handler into the middleware.

#### Error handling middleware

It's the same situation as the net/http example. The cody is pretty much copy-pasted from there with minor changes.

chi itself doesn't have an error handling middleware readily available, so in order to make that work, we need to shove the errors into the context, and make sure the wrapping middleware also has access to the new one, so gotta do the ol' `*r = *r.WithContext(ctx)` dance.

#### Context up and down

As above, kinda clunky, but easily solvable.

#### Unit tests

Standard `httptest` can be used. Absolutely fine to use, same as every other one.

#### Ecosystem

chi has a few [middlewares built in](https://go-chi.io/#/pages/middleware) ([code here](https://github.com/go-chi/chi/tree/master/middleware)).

There doesn't seem to be any contrib repository, but also they wouldn't be chi specific, because a middleware is just a `func(http.Handler) http.Handler`, which is standard library.
