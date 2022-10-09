# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## gin implementation
### General considerations
#### paths with multiple http verbs
Gin is... kind of weird to get started with. Compared to [echo](https://github.com/suborbital/framework-muxer-showdown/tree/echo), the only two ways to define a route with multiple verbs are
* either multiple single declarations, like `gin.GET` and `gin.OPTIONS` for the same route, or
* `gin.Any()`, which will enable the route for all verbs, and then put a configured middleware on that one specific route

If I choose the single declarations for the route, then any request to a verb that's not supported gets a 404 instead of a 405. Whether that's something you want or not is a different question, but the any route + middleware can be used to return a correct empty 405 response.

I'm unsure which solution I like less between them.

#### graceful shutdown
By default, gin is supposed to be used as
```go
r := gin.Default()
r.GET("/path", handlerFunc)

r.Run(":3000")
```
Nice and easy, except there's no way to do graceful shutdown. There's no `r.Stop()` or `r.Shutdown(ctx)` or anything. Instead, [the readme suggests (towards the end)](https://github.com/gin-gonic/gin#graceful-shutdown-or-restart) to use either a 3rd party library, or to use standard library `http.Server` and the gin engine as a handler in the server.
```go
router := gin.Default()

// Health endpoint
router.GET("/path", handlerFunc)

server := &http.Server{
	Addr:    ":9000",
	Handler: router.Handler(),
}
```
This way when it's time to stop the service I can `server.Shutdown(ctx)` and be graceful about it rather than terminating the service and hoping that `router.Run()` cleans up after itself.

**EDIT:** I checked, it does not clean up after itself. `router.Run()` essentially does an `http.ListenAndServe(address, router.Handler())`, which is essentially the following anyways:
```go
http.Server{
	address: address,
	handler: router.Handler(),
}
```
This blocks, so if the service is terminated, every active connection is simply yanked.

There is no functional difference between attaching gin to an http server and starting that, so we can gracefully shut down, or running gin directly from a server running standpoint.

It works, the lack of convenience in this part is just mildly annoying.

#### nice debug output
This is with debug mode on. Debug mode should be turned off when running in production, it will be significantly less noisy.

![colorized gin terminal](assets/gin%20terminal.jpg)

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
