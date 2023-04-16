# framework-muxer-showdown

A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

- `make test`: runs `go test ./...`.
- `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
- `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
- `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## httptreemux implementation

### General Considerations

See the implementation here: https://github.com/javorszky/framework-muxer-showdown/tree/httptreemux

There are two ways of working with this - httptreemux.New() or httptreemux.NewContextMux(). The later uses http.Handler and http.HandlerFunc. The former uses its own function signature with a third map[string]string for params.

### Details of criteria

#### Context type

This is an embedded context into the http request.

#### Can I use it as http.Handler
Yes
```go
r := httptreemux.NewContextMux()

var treemuxIsHandler http.Handler

treemuxIsHandler = r

r.GET("/router-is-handler", treemuxIsHandler.ServeHTTP)
```
#### Standard library handling

This comes out of the box because we are using the NewContextMux.

#### Accessing raw Request and ResponseWriter

This comes out of the box because it uses the http.Handler and http.HandlerFunc.

#### Websocket

Works the same as the net/http implementation, or the chi implementation.

#### Path specificity

This is a kinda.

* on one hand, having a catch-all, and a longer static route works. The static route will match every time first, if it can
* on the other hand, having a single with `/spec`, and the catch-all at `/spec/*stuff`, a request to `GET /spec/` will match the single rather than the catch-all handler

#### Path variables

Yep, works, even with the standard handler muxer, because of a helper function: `httptreemux.ContextParams(r.Context())`, so it just works.

#### Grouping

Grouping works, and you can add a bunch of middlewares to each group as well.

#### Overlaps

This works as expected.

#### General middleware

I like the way middleware handling works here. As there are two different types of routers, there are also two different ways of handling middlewares when it comes to function signatures. The implementation is with the http handlers.

You can add middlewares to both the global handler, and also to each of the groups. They just work.

Layering test
```go
r := httptreemux.NewContextMux()
r.UseHandler(handlers.MidOne(l))
r.UseHandler(handlers.MidTwo(l))

r.GET("/layer",
    handlers.MidThree(l)(
        handlers.MidFour(l)(
            handlers.StandardHandlerFunc(),
        ),
    ).ServeHTTP,
)
```
```shell
{"level":"info","component":"app","mid":"one","message":"red one reporting in"}
{"level":"info","component":"app","mid":"two","message":"red two reporting in"}
{"level":"info","component":"app","mid":"three","message":"red three reporting in"}
{"level":"info","component":"app","mid":"four","message":"red four reporting in"}
```

#### Error handling middleware

No special handling for error handling middleware, so the solution is much like net/http. That said a copy-paste for both the handlers and the middleware was enough to make it work.

#### Context up and down

Same as net/http, need to do the *r = *r.WithContext dance. Otherwise works.

#### Unit tests

Standard httptest works. Both with the handlers themselves, without involving the muxer, or with the muxer itself. In both cases we pass the test request and test recorder to the handler func, the `ServeHTTP` method both on the actual router, and the individual handlers.

The only thing we need to take care of is that while we can use the `router.UseHandler` to specify global middlewares and then use the `router.ServeHTTP` to test handlers with paths, with individual handlers we do need to wrap them up.

#### Ecosystem

This is a `Some`, because it relies heavily on standard library solutions, and there isn't a lot of community around it, because it's non-specific. There are a bunch of other frameworks built on top of the muxer, but they aren't listed anywhere on the readme.

#### Performance

httptreemux has router configurable panic handler, so that will be present for `/smol-perf` too.

##### /performance

* [perftests/run1.log](perftests/run1.log):  Requests/sec:	22492.1095
* [perftests/run2.log](perftests/run2.log):  Requests/sec:	22290.5159
* [perftests/run3.log](perftests/run3.log):  Requests/sec:	22238.6294
* [perftests/run4.log](perftests/run4.log):  Requests/sec:	22591.5592

##### /smol-perf

* [perftests/smol-run1.log](perftests/smol-run1.log):  Requests/sec:	14754.3833
* [perftests/smol-run2.log](perftests/smol-run2.log):  Requests/sec:	21348.5404
* [perftests/smol-run3.log](perftests/smol-run3.log):  Requests/sec:	22246.9219
* [perftests/smol-run4.log](perftests/smol-run4.log):  Requests/sec:	15226.8435
