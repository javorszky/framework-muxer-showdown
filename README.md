# framework-muxer-showdown

A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

- `make test`: runs `go test ./...`.
- `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
- `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
- `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## Implementations and tests

## httptreemux

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

#### Path specificity

#### Path variables

#### Grouping

#### Overlaps

#### General middleware

#### Error handling middleware

#### Context up and down

#### Unit tests

#### Ecosystem
