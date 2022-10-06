# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

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

Error handling is generally the weakest point as we'd need to respond from the end handlers directly, because none of the handlers have return arguments that we could catch.

### Details of criteria

#### Context type

#### Standard library handling

#### Accessing raw Request and ResponseWriter

#### Websocket

#### Path specificity

#### Path variables

#### Grouping

#### Overlaps

#### General middleware

#### Error handling middleware

#### Context up and down

#### Unit tests

#### Ecosystem
