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

Yep, there's an `echo.WrapHandler(http.Handler)` that's readily availably. Super easy to use.

#### Accessing raw Request and ResponseWriter

It's easy to do, and makes using websockets clean.
```go
func SomeHandler(c echo.Context) error {
	req := c.Request()
	writer := c.Response()
	
	// the rest of the handler
}
```
See the [documentation on using websockets](https://echo.labstack.com/cookbook/websocket/) for an example usage).

#### Websocket

As above, it's straightforward. Echo doesn't provide its own implementation of it, but rather recommends you use one of the other websocket libraries, whether that's standard library in the /x/ namespace, or gobwas, or gorilla's websocket, it doesn't matter.

As long as you have access to the request and responsewriter, you can upgrade the incoming GET request.

Keep in mind that if you use the standard library, the request MUST have an Origin header with a value that's a parseable `url.URL`.

An implementation is in [the websocket handler](handlers/ws.go).

#### Path specificity

#### Path variables

#### Grouping

#### Overlaps

Echo beautifully handles overlaps, I can have a route that's dynamic, and also static on the same prefix, and both work as expected. See the [implementation of overlaps](handlers/overlaps.go) and the declaration of the [routes in app.go](app/app.go).

More surprisingly it also handles the case where the dynamic part is completely missing, though a new catch-all declaration with `/prefix/` is needed.

I am very pleased with it!

#### General middleware

Super easy to use, works a lot similar to the [middleware situation in the net/http](https://github.com/suborbital/framework-muxer-showdown/tree/net/http#middlewares-easy) implementation.

This is essentially what it looks like:

```go
func ExampleMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// do things before calling the wrapped handler
			logger.Info().Str("middleware", "MidTwo").Msg("goodbye!")

			err := next(c)
			
			// do things after calling the wrapped handler
			
			return err
		}
	}
}
```
It's an onion type, and we'd plug this in like so:
```go
func New() {
	e := echo.New()
	
	// have to call that func to return the middleware type.
	e.Use(mid.ExampleMiddleware())
	e.Use(mid.OtherMiddleware())
	
	// or for individual routes
	e.Get("/some-path", handlers.SomePathHandler(), mid.PathMidOne(), mid.PathMidTwo())
}
```
Bear in mind that when handling the request, the latter middlewares in `e.Use` get called first, but the latter middlewares in a path declaration get called last, so the above would be:
```go
Req -> mid.Other -> mid.Example -> mid.PathMidOne -> mid.PathMidTwo -> handlers.SomePathHandler
```

#### Error handling middleware

The exported `echo.HTTPErrorHandler` can be reassigned to be a custom error handler, and for all misc errors we don't want to handle, we can forward it to the also exported `echo.DefaultHTTPErrorHandler`. It's fairly easy.

See the [ErrorHandler exported func in middlewares](handlers/middlewares.go).

Echo handlers also return an error, so that makes working with it kind of a breeze too!

#### Context up and down

#### Unit tests

This one is kind of awkward, but in the end it's a `Yes`, rather than a `Kinda`. We can still use the `httptest` package for a new request and new recorder. There are two things to watch out for with echo:

1. Because the handler funcs expect echo's own context interface, we need to create a new echo, and then a new context based on the httptest request and recorders. That's easy to do, but needs to be done. See the [health endpoint test](handlers/health_test.go)
2. If we're testing endpoints where we need to rely on some ✨_M A G I C_✨ that echo gives us, we need to first add the handler to a new echo instance, and test the `echo.ServeHTTP` instead of the handler itself. The end result is the same, because that's the only handler that _should_ be on the echo instance anyways, but it also has all the other scaffolding like [the custom error handler](handlers/errors_test.go) or [setting up binding for path variables](handlers/pathvars_test.go)

#### Context up-down

`echo.Context` has `Get(key string)` and `Set(key string, value interface{})` methods that make it easy to pass information easily between layers of middlewares and handlers. We also don't need to do magic by assigning pointers to pointers, so it's a whole lot easier.

See the `ContextUpDown` [middleware implementation](handlers/middlewares.go) and the `UpDownHandler` [handler function](handlers/contextupdown.go) on how they work. The output to look for is in the terminal, the request / response is not involved in this.

#### Ecosystem
