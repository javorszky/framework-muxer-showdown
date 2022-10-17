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

See the implementation here: https://github.com/suborbital/framework-muxer-showdown/tree/httptreemux

There are two ways of working with this - httptreemux.New() or httptreemux.NewContextMux(). The later uses http.Handler and http.HandlerFunc. The former uses its own function signature with a third map[string]string for params.

### Details of criteria

#### Context type

This is an embedded context into the http request.

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

#### Error handling middleware

No special handling for error handling middleware, so the solution is much like net/http. That said a copy-paste for both the handlers and the middleware was enough to make it work.

#### Context up and down

#### Unit tests

#### Ecosystem
