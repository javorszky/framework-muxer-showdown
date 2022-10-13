# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## fiber implementation

### General Considerations

Hnng, the configuration option! DAMN! I love it!

Fiber itself has `.Listen()` and `.Shutdown()` methods, so we don't need to involve `http.Server` directly.

Has configurable ErrorHandler!

Also has an option to limit requests to be GET only.

The `Prefork` option is super interesting!

`StrictRouting` needs to be set to true, otherwise `/spec` and `/spec/` are treated as the same.

Dependency graph is small! The entire `go.mod` file of fiber is this:
```go
module github.com/gofiber/fiber/v2

go 1.19

require (
	github.com/valyala/fasthttp v1.40.0
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
)
```

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
