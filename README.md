# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## julienschmidt/httprouter implementation

### General Considerations

* This is just a router. Ardan Labs's Service uses it.
* It has a stated aim of routes matching exactly one, or zero routes, which means it will probably fail the specificity and overlap tests, but we'll see
* Doesn't have its own start / stop, so it goes on the `&http.Server{Handler: routerInstance}`. The upside is that it's easy to start / stop in a standard way.

### Details of criteria

#### Context type

Embedded into the standard `*http.Request` type.

#### Standard library handling

There are convenient adapters to deal with it. The parameters are then stored in the request context.

#### Accessing raw Request and ResponseWriter

Its own handler type has a signature of `func(http.ResponseWriter, *http.Request, httprouter.Params) {}`. The raw request and response writers are just... there.

#### Websocket

#### Path specificity

#### Path variables

Yep, it supports it. In the `httprouter.Handle` signature, the third parameter is the parameters in the path, but query parameters aren't contained there.

#### Grouping

#### Overlaps

#### General middleware

Eh, it's ... easy. Sort of. It's missing a lot of the convenience methods like `.Use` as the others have, but you can wrap them inside each other.

The router has a `.Handle` method that takes a standard http handler interface and that provides a wrapper around that so from the outside it looks like an `httprouter.Handle`. I wish that was exported, but at the same time it's about 8 lines, so moved it to the `handlers.Wrap` method.

#### Error handling middleware

#### Context up and down

#### Unit tests

#### Ecosystem
