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

Mostly all right. The way it works is not really comfortable, but it gets the job done, mostly. No standard library context.Context, though there's a flag on the `gin` router that can be set to enable the timeout / deadline / cancel / done methods on it.

## Template for readme

### General Considerations

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
