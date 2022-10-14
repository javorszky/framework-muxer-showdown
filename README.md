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
* Has good support for CORS things: https://github.com/julienschmidt/httprouter#automatic-options-responses-and-cors

### Details of criteria

#### Context type

Embedded into the standard `*http.Request` type.

#### Standard library handling

There are convenient adapters to deal with it. The parameters are then stored in the request context.

#### Accessing raw Request and ResponseWriter

Its own handler type has a signature of `func(http.ResponseWriter, *http.Request, httprouter.Params) {}`. The raw request and response writers are just... there.

#### Websocket

No built in support, but the standard library just works, mostly because it's an `http.Handler`.

#### Path specificity

There is a concept of wildcard parameters, which works, but you can't also declare a handler for a route that is more specific. Some solutions allow you to declare them, just won't work.

This test fails with the following panic given the route declarations:
```go
	// Path specificity
	r.GET("/spec/*stuff", handlers.Everyone())
	r.GET("/spec", handlers.Single())
	r.GET("/spec/long/url/here", handlers.Long())
```
```shell
panic: '/long/url/here' in new path '/spec/long/url/here' conflicts with existing wildcard '/*stuff' in existing prefix '/spec/*stuff'

goroutine 1 [running]:
<rest of the stacktrace>
```

#### Path variables

Yep, it supports it. In the `httprouter.Handle` signature, the third parameter is the parameters in the path, but query parameters aren't contained there.

#### Grouping
This is ... weird, because it's not even supported as much as net/http. There's an issue: https://github.com/julienschmidt/httprouter/pull/89, in 2016 there was a promise of a "new version", but current version does not support it.

There's a separate middleware? Or a module here: https://github.com/omgnuts/go-subware

But it's... an eeeeehhhh, I'd rather don't want to deal with this.

I tried to mount a sub router to a route, and it does not work.

#### Overlaps

Also does not work, and exits with a panic given the following declaration:
```go
	// Overlaps
	r.GET("/overlap/:one", handlers.OverlapDynamic())
	r.GET("/overlap/kansas", handlers.OverlapSpecific())
	r.GET("/overlap/", handlers.OverlapEveryone())
```

```shell
panic: 'kansas' in new path '/overlap/kansas' conflicts with existing wildcard ':one' in existing prefix '/overlap/:one'

goroutine 1 [running]:
```

#### General middleware

Eh, it's ... easy. Sort of. It's missing a lot of the convenience methods like `.Use` as the others have, but you can wrap them inside each other.

The router has a `.Handle` method that takes a standard http handler interface and that provides a wrapper around that so from the outside it looks like an `httprouter.Handle`. I wish that was exported, but at the same time it's about 8 lines, so moved it to the `handlers.Wrap` method.

#### Error handling middleware

Same as the net/http solution. Copy pasted the code from there, minimal modification was needed.

#### Context up and down

#### Unit tests

#### Ecosystem
