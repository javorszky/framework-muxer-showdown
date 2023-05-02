# framework-muxer-showdown

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
Gin is... kind of weird to get started with. Compared to [echo](https://github.com/javorszky/framework-muxer-showdown/tree/echo), the only two ways to define a route with multiple verbs are
* either multiple single declarations, like `gin.GET` and `gin.OPTIONS` for the same route, or
* `gin.Any()`, which will enable the route for all verbs, and then put a configured middleware on that one specific route

If I choose the single declarations for the route, then any request to a verb that's not supported gets a 404 instead of a 405. Whether that's something you want or not is a different question, but the `any` route + middleware can be used to return a correct empty 405 response.

I'm unsure which solution I like less between them.

⚠️ **EDIT**: found how to do it the gin way. Each engine / group has an option to `HandleMethodNotFound`, which is set to `false` by default. That means if a given verb+route is requested, but that doesn't have a handler, it gets a 404. If that setting is set to `true`, it will check whether other verbs for the same path exist, and then it sends it to a `NoMethod` handler.

The `NoMethod` handler can also be customized, which I've done to remove the body from the response.


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

**EDIT:** I checked, it does not clean up after itself. `router.Run()` essentially does an `http.ListenAndServe(address, router.Handler())`, which is essentially the following anyway:
```go
http.Server{
	address: address,
	handler: router.Handler(),
}
```
This blocks, so if the service is terminated, every active connection is simply yanked.

There is no functional difference between attaching gin to an HTTP server and starting that, so we can gracefully shut down, or running gin directly from a server running standpoint.

It works, the lack of convenience in this part is just mildly annoying.

#### nice debug output
This is with debug mode on. Debug mode should be turned off when running in production, it will be significantly less noisy.

![colorized gin terminal](assets/gin%20terminal.jpg)

### Details of criteria
#### Context type

It's their own context, and not even an interface, like in the case of echo, but an actual concrete struct.

Granted it has a LOT of moving parts and capabilities, from storing the original request and ResponseWriter, to shorthands for sending data back to the client, to aborting the request due to an error, to setting / getting values up and down the chain, and storing the chain of handlers for a given request.

#### Can I use it as http.Handler?
Yes
```go
router := gin.New()

var ginIsHandler http.Handler

ginIsHandler = router

router.Any("/router-is-handler", gin.WrapH(ginIsHandler))
```

#### Standard library handling

Yes, gin provides `WrapF(f http.HandlerFunc)` and `WrapH(h http.Handler)` functions that turn both into a `gin.HandlerFunc(c *gin.Context) {}`, so this is straightforward.

#### Accessing raw Request and ResponseWriter

Yes, they're on `c.Writer` and `c.Request` exported properties.

#### Websocket

Does not have anything built in, so upgrading the GET request to a WS is up to us. Gobwas or standard library /x/websocket is perfectly fine for the job.

As we can access the raw request and ResponseWriter, it's essentially the same solution as echo's or net/http's.

#### Path specificity

Gin isn't as powerful as [echo](https://github.com/javorszky/framework-muxer-showdown/tree/echo) in this. It handles the single, and everyone else cases, but having these three declarations at the same time causes a panic:
```go
// Path specificity
router.GET("/spec", handlers.Single())
router.GET("/spec/*thing", handlers.Everyone())
router.GET("/spec/long/url/here", handlers.LongRoute())
```
```shell
[GIN-debug] GET    /spec                     --> github.com/javorszky/framework-muxer-showdown/handlers.Single.func1 (3 handlers)
[GIN-debug] GET    /spec/*thing              --> github.com/javorszky/framework-muxer-showdown/handlers.Everyone.func1 (3 handlers)
[GIN-debug] GET    /spec/long/url/here       --> github.com/javorszky/framework-muxer-showdown/handlers.LongRoute.func1 (3 handlers)
panic: '/long/url/here' in new path '/spec/long/url/here' conflicts with existing wildcard '/*thing' in existing prefix '/spec/*thing'

goroutine 1 [running]:
<stacktrace here>
```
The two ways out of this are

##### Define one catch-all wildcard route

And then inspect what the wildcard is, create a new gin context, and pass it on to a handler. Potentially passing in a sub-muxer gin engine with all the specific routes.

##### Use a custom 404 handler middleware

Define your catch-all, and the long specific routes, like so:
```go
router.GET("/spec", handlers.Single())
router.GET("/spec/",  handlers.Everyone())
router.GET("/spec/long/url/here", handlers.LongRoute()) // this one doesn't work with the above
```
This way `/spec` will match the single, `/spec/` will match the everyone else, `/spec/long/url/here` will match the `LongRoute`, but `/spec/somewhere/warm` would get a 404. This is where the custom 404 handler comes in: each of those can be captured and rerouted to the EveryoneElse handler.

However, that means that the global error handler needs to be modified, and at that point that could become a really really big mess of spaghetti code and god function.

#### Path variables

Unsurprising, and works fairly well. There's not much to write home about.

```go
// /pathvars/:one/metrics/:two

func (c *gin.Context) {
	firstParam  := c.Param("one")
	secondParam := c.Param("two")
	// the rest of the owl
}
```

#### Grouping

Implemented two different groups: one for the path specificity, the other for the actual group endpoint. They just work, 10/10 no notes.

#### Overlaps

Also just kinda work, the gin documentation (readme) also calls it out.

#### General middleware

Middlewares and handlers are the same signature, the only difference is that in a middleware you can call `c.Next()`, which will call the next handler in the chain.

Layering test with the following:
```go
router := gin.New()
router.Use(handlers.MidOne(logger))
router.Use(handlers.MidTwo(logger))
router.GET("/layer",
	handlers.MidThree(logger),
	handlers.MidFour(logger),
	gin.WrapH(handlers.StandardHandler()),
)
```
```shell
{"level":"info","component":"app","mid":"one","message":"red one reporting in"}
{"level":"info","component":"app","mid":"two","message":"red two reporting in"}
{"level":"info","component":"app","mid":"three","message":"red three reporting in"}
{"level":"info","component":"app","mid":"four","message":"red four reporting in"}
```

#### Error handling middleware

I put this down as `Kinda`. It's a lot more clunky than I would like. Gin has the concept of errors on context, you can grab the errors from the context, but it's going to be a slice, and then you need to sort through all of them and figure out which one is the most important to return to the client.

On top of that there isn't a built in error handler, or a default error handler. The two that are present are the notfound and nomethod handlers.

It can be done, but it's finnicky, and... not comfortable, for lack of a better description.

#### Context up and down

Yeah, it works. `c.Set(key string, value any)` and `c.Get(key string)` work as expected. See the [ctx updown example](handlers/ctxupdown.go), the other part of it is in the [middlewares file](handlers/middlewares.go).

#### Unit tests

Unit testing is also done with the `httptest` standard library. As opposed to echo, there's no way to create a `gin.Context` from the request and response writer, we gotta put together the entire muxer and configure it, like in the [health test](handlers/health_test.go) file, and then test the `.ServeHTTP(w, r)` method of it.

Other than that it's fairly straightforward, and this helps us test middlewares and handlers in isolation.

#### Ecosystem

There are plenty of middlewares and whatnot.

* examples live here: https://github.com/gin-gonic/examples
* 3rd party middlewares: https://github.com/gin-contrib (includes logging with rs/zerolog, a cors middleware, authz, and requestID)

#### Performance

Gin doesn't have router configurable panic or error handlers.

##### /performance

* [perftests/run1.log](perftests/run1.log):  Requests/sec:	30864.4914
* [perftests/run2.log](perftests/run2.log):  Requests/sec:	32513.2556
* [perftests/run3.log](perftests/run3.log):  Requests/sec:	32970.2756
* [perftests/run4.log](perftests/run4.log):  Requests/sec:	32800.8955

##### /smol-perf

* [perftests/smol-run1.log](perftests/smol-run1.log):  Requests/sec:	23343.5163
* [perftests/smol-run2.log](perftests/smol-run2.log):  Requests/sec:	15332.3898
* [perftests/smol-run3.log](perftests/smol-run3.log):  Requests/sec:	21377.7060
* [perftests/smol-run4.log](perftests/smol-run4.log):  Requests/sec:	21775.2100
