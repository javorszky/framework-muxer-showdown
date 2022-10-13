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

It's a custom `*fiber.Ctx` type. Works similarly to gin in that it has `c.Next()`.

#### Standard library handling

This is super nice! Most other solutions either work like that (gin, chi), or have the wrappers built into the router itself (echo). Fiber's approach is to delegate this to a middleware entirely, so if you don't need it, you don't have the additional code overhead, but if you do need it, it's ready.

Better yet, the middleware [is on their official GitHub org](https://github.com/gofiber/adaptor), and does conversion in both ways:

http handler / handlerfunc <-> fiber handler
http middleware <-> fiber middleware
fiber.App -> http handlerfunc

This is super nice! Major props and thumbs up!

#### Accessing raw Request and ResponseWriter

I was trying to figure out whether this should be a `Yes` or `Kinda`, but ultimately went with `Yes`, because if you need the raw `http.Request` or the `http.ResponseWriter`, you're probably going to be happy using the `http.Handler` wrapper middleware.

There are two layers of wraps, one is fiber, the other is fasthttp within fiber, and fasthttp, from the looks of it, is not super forthcoming with letting the request / responsewriter be accessible, though this might change once I look into fasthttp proper.

#### Websocket

Fiber has its own websocket implementation, which I suppose is a wrapper. It lives in a completely different repository: [gofiber/websocket](https://github.com/gofiber/websocket), but actually using it is fairly straightforward.

There's no need to involve other libraries, like the std library websocket, gobwas, or gorilla's websocket.

#### Path specificity

Handles path specificity tests. Wildcard and a more specific route that overlaps with it is good to go.

#### Path variables

#### Grouping

#### Overlaps

#### General middleware

General middleware looks the same as a handler itself. If it's a middleware, there will be a `c.Next()` call that we can do to go down the chain.

#### Error handling middleware

There is a special type of error handling type, `fiber.ErrorHandler`. It takes a ctx and the error, and returns an error.

Fiber has a [default error handler](https://github.com/gofiber/fiber/blob/6a5fc64eddaa81a7fb65c94b8dcfd9a2caac2e78/app.go#L452-L461) which either returns the code and message from an internal `fiber.Error`, or returns a standard internal server error to the client.

There's one place to handle all errors, so in this regard it's similar to echo. Redefining and adjusting the error handler is very easy. One of `fiber.Config`'s properties is the error handler.

#### Context up and down

#### Unit tests

#### Ecosystem
