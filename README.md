# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## echo implementation

### General Considerations

See the tree here: https://github.com/suborbital/framework-muxer-showdown/tree/echo

API is nice, has a match-any method, a match-some method, and a match-single method handler.

### Details of criteria

#### Context type

It's a custom thing. :( Even though the standard `context.Context` is an interface. Technically, because both the std context, and echo's context are interfaces, we _could_ merge them both into a custom context and use that, but that's a bit annoying.

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
