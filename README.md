# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## Implementations and tests
### 1. net/http

Pretty basic, it gets us to about 80%. Only dependency we really need is when we want to deal with websockets and don't want to use the /x/ standard library.

Major downside is path variables. They _can_ be done, but that's a lot of extra code.

See the tree here: https://github.com/suborbital/framework-muxer-showdown/tree/net/http

### 2. echo implementation

Very clean, very simple to use, supports everything we really need. Only downside is the lack of standard library `context.Context`, but the timeouts can be configured on the echo instance before startup.

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/echo

### 3. gin

Mostly all right. The way it works is not really comfortable, but it gets the job done. No standard library context.Context, though there's a flag on the `gin` router that can be set to enable the timeout / deadline / cancel / done methods on it.

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/gin
### 4. chi

chi is mostly similar to the standard library net/http implementation with its very very standard signatures, with the added benefit of url params and routing.

Between net/http and chi, chi wins.
Between chi and gin, chi wins, because gin can't do a routing we need.
Between chi and echo though, echo wins because of significantly easier error handling.

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/chi
### 5. fiber

I **really** like fiber, despite the fact that it has a custom ctx (this seems to be a common theme), despite the fact we can't easily access the http request and response writers, and despite the fact that unit testing doesn't use the NewRecorder, and despite the weird ordering need to make the overlap happen.

It makes up for all of those by providing convenience methods and middlewares that just kind of make sense.

It has a really robust configuration option, and grouping and middlewares are excellent, and clear.

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/fiber

### 6. httprouter

This is just a router rather than a web framework. It forms the basis of Gin, and also Ardan Labs's Service starting boilerplate implementation.

It has a bunch of decisions in it, like each request can only match one or none routes, which means the path specificity and overlap tests fail on our end. Whether those are good decisions or not depends on the use case, but I'd wager that it's inconvenient for us.

Grouping is also very problematic. I suppose Gin fixed some of the issues and made it more convenient to work with.

I can't recommend we use it.

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/httprouter

### 7. httptreemux

Just a router, super similar to httprouter. There's only one test that really fails, otherwise everything else is super nice.

I particularly like that I can create a muxer that's either the standard library handlers, or its own signature with the params as a third one.

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/httptreemux

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

#### Overlaps

#### General middleware

#### Error handling middleware

#### Context up and down

#### Unit tests

#### Ecosystem
