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

#### Websocket

#### Path specificity

#### Path variables

#### Grouping

#### Overlaps

#### General middleware

Super easy to use, works a lot similar to the [middleware situation in the net/http](https://github.com/suborbital/framework-muxer-showdown/tree/net/http#middlewares-easy) implementation.

This is essentially what it looks like:

```go
package mid

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
	
	// gotta call that func to return the middleware type.
	e.Use(mid.ExampleMiddleware())
	e.Use(mid.OtherMiddleware())
	
	// or for individual routes
	e.Get("/somepath", handlers.SomePathHandler(), mid.PathMidOne(), mid.PathMidTwo())
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

#### Ecosystem
