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

Observations, in no particular order:

* Hnng, the configuration option! DAMN! I love it!
* Fiber itself has `.Listen()` and `.Shutdown()` methods, so we don't need to involve `http.Server` directly.
* Has configurable ErrorHandler!
* Also has an option to limit requests to be GET only.
* The `Prefork` option is super interesting!
* `StrictRouting` needs to be set to true, otherwise `/spec` and `/spec/` are treated as the same.
* Mercedes-Benz is a GitHub sponsor for them (yes, the car company)

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
Though this does not include the recover middleware, nor the adaptor middleware.

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

So far [fiber's path variables](https://docs.gofiber.io/guide/routing#parameters) are the most robust, including a regex, as well as a bunch of other constraints, trumping chi's regex constraints!

#### Grouping

Very nice, it does all the things it needs to do like group level middlewares, etc, but you can also define a handler for the prefix, so if you have a group of `/v1/*`, you can then define `/hello` on the group, and the full path becomes `/v1/hello`. When you define a handler on the group though, going to `/v1` or `/v1` will also work.

You can also mount a group to a group. The [documentation on groups](https://docs.gofiber.io/guide/grouping) is pretty straightforward too.

#### Overlaps

So this is weird, because it depends on the order you declare the routes. So this one works, the kansas one goes to handler, everything else goes to the dynamic handlers:
```go
f.Get("/overlap/kansas", handler)
f.Get("/overlap/:one", dynamicHandler)
```
However, this one doesn't work as expected, every request, including /overlap/kansas goes to the dynamic handler. 
```go
f.Get("/overlap/:one", dynamicHandler)
f.Get("/overlap/kansas", handler)
```

Also tried it the group way, but did not help, still depended on the order. This might be a bug in fasthttp.

#### General middleware

General middleware looks the same as a handler itself. If it's a middleware, there will be a `c.Next()` call that we can do to go down the chain.

#### Error handling middleware

There is a special type of error handling type, `fiber.ErrorHandler`. It takes a ctx and the error, and returns an error.

Fiber has a [default error handler](https://github.com/gofiber/fiber/blob/6a5fc64eddaa81a7fb65c94b8dcfd9a2caac2e78/app.go#L452-L461) which either returns the code and message from an internal `fiber.Error`, or returns a standard internal server error to the client.

There's one place to handle all errors, so in this regard it's similar to echo. Redefining and adjusting the error handler is very easy. One of `fiber.Config`'s properties is the error handler.

#### Context up and down

Yep, this works, but unlike the others, there's a `.Locals(key, ...values)` method. If you only supply the key, it reads the value stored. If you supply both a key and a value, that sets the value.

#### Unit tests

Unit testing is also done via http test. The only minor difference between this one and the others is that instead of using a `NewRecorder` as a writer, fiber has a `.Test` method to make working with tests a little bit easier.

I found no easy way to use the NewRecorder, as fiber and fasthttp seem to abstract those behind the `.Test` method.

See the [errors_test.go](handlers/errors_test.go) file for more details.

#### Ecosystem

[Fiber's GitHub org](https://github.com/gofiber) has a bunch of useful repositories, middlewares. But other than that I've not seen extensive support / list of community bits. They do have a [discord however](https://docs.gofiber.io/extra/faq#does-fiber-have-a-community-chat).

I've upgraded it to `Rich` because of this: https://github.com/gofiber/awesome-fiber.

#### Performance

Fiber has a router configurable error handler, so that will also wrap the `/smol-perf` route.
##### With standard json.Marshal/Unmarshaler

###### /performance

[perftests/stdjson/run1.log](perftests/stdjson/run1.log):  Requests/sec:	64937.3429
[perftests/stdjson/run2.log](perftests/stdjson/run2.log):  Requests/sec:	54957.3747
[perftests/stdjson/run3.log](perftests/stdjson/run3.log):  Requests/sec:	64256.2797
[perftests/stdjson/run4.log](perftests/stdjson/run4.log):  Requests/sec:	58466.2158

###### /smol-perf
[perftests/stdjson/smol-run1.log](perftests/stdjson/smol-run1.log):  Requests/sec:	46897.7520
[perftests/stdjson/smol-run2.log](perftests/stdjson/smol-run2.log):  Requests/sec:	50371.1567
[perftests/stdjson/smol-run3.log](perftests/stdjson/smol-run3.log):  Requests/sec:	31360.7054
[perftests/stdjson/smol-run4.log](perftests/stdjson/smol-run4.log):  Requests/sec:	52037.3624

#### bytedance/sonic

https://github.com/bytedance/sonic

##### /performance

[perftests/sonic/run1.log](perftests/sonic/run1.log):  Requests/sec:	65296.1769
[perftests/sonic/run2.log](perftests/sonic/run2.log):  Requests/sec:	64270.4320
[perftests/sonic/run3.log](perftests/sonic/run3.log):  Requests/sec:	63664.6592
[perftests/sonic/run4.log](perftests/sonic/run4.log):  Requests/sec:	63249.9146

##### /smol-perf

[perftests/sonic/smol-run1.log](perftests/sonic/smol-run1.log):  Requests/sec:	32491.4713
[perftests/sonic/smol-run2.log](perftests/sonic/smol-run2.log):  Requests/sec:	34704.7958
[perftests/sonic/smol-run3.log](perftests/sonic/smol-run3.log):  Requests/sec:	29288.5351
[perftests/sonic/smol-run4.log](perftests/sonic/smol-run4.log):  Requests/sec:	47309.3414
