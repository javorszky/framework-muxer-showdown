# framework-muxer-showdown

## Scaffolding

Makefile has a docker build that produces a utility image with `gci`, `golangci-lint`, and `mockery` in it. The docker compose file then uses that image to run the various tools, so you don't need to have any of them installed locally.

### Commands

* `make test`: runs `go test ./...`.
* `make lint`: runs golangci-lint with the configs from the [.golangci.yaml](.golangci.yaml) file.
* `make lintfix`: runs gci on all `*.go` files recursively minus the `vendor` directory.
* `make mocks`: runs mockery to generate mocked interfaces in all go files recursively minus the `vendor` directory with config from the [.mockery.yaml](.mockery.yaml) file.

## fasthttp implementation

### General Considerations

There are two things we can refer to as "fasthttp". One of them is merely the handler, which you can find at [https://github.com/valyala/fasthttp](https://github.com/valyala/fasthttp), and the other one is [https://github.com/fasthttp/router](https://github.com/fasthttp/router), a router based on it even though the readme says it's based on `julienschmidt/httprouter`.

The test itself is about the `fasthttp/router`, because the naked handler has the following under here:

<details>
<summary>Reasons we're not using `valyala/fasthttp`</summary>

---

Its main focus is on performance and it compares itself to net/http. It also starts with this note:

> ## fasthttp might not be for you!
>
> ---
> fasthttp was design for some high performance edge cases. Unless your server/client needs to handle thousands of small to medium requests per seconds and needs a consistent low millisecond response time fasthttp might not be for you. For most cases net/http is much better as it's easier to use and can handle more cases. For most cases you won't even notice the performance difference.

I feel like in our case we do care about these fast responses, so it might be for us.

That said, fasthttp also has this bit in their readme:

> * Fasthttp doesn't provide [ServeMux](https://golang.org/pkg/net/http/#ServeMux), but there are more powerful third-party routers and web frameworks with fasthttp support:
>   * [fasthttp-routing](https://github.com/qiangxue/fasthttp-routing)
>   * [router](https://github.com/fasthttp/router)
>   * [lu](https://github.com/vincentLiuxiang/lu)
>   * [atreugo](https://github.com/savsgio/atreugo)
>   * [Fiber](https://github.com/gofiber/fiber)
>   * [Gearbox](https://github.com/gogearbox/gearbox)
>
> Net/http code with simple ServeMux is trivially converted to fasthttp code:
>
> ```go
> // net/http code
>
> m := &http.ServeMux{}
> m.HandleFunc("/foo", fooHandlerFunc)
> m.HandleFunc("/bar", barHandlerFunc)
> m.Handle("/baz", bazHandler)
>
> http.ListenAndServe(":80", m)
> ```
>
> ```go
> // the corresponding fasthttp code
> m := func(ctx *fasthttp.RequestCtx) {
> 	switch string(ctx.Path()) {
> 	case "/foo":
> 		fooHandlerFunc(ctx)
> 	case "/bar":
> 		barHandlerFunc(ctx)
> 	case "/baz":
> 		bazHandler.HandlerFunc(ctx)
> 	default:
> 		ctx.Error("not found", fasthttp.StatusNotFound)
> 	}
> }
>
> fasthttp.ListenAndServe(":80", m)
> ```

All of the above, while undoubtedly amazing when it comes to handling requests, is going to make it really, really difficult to write actual code that we need as a business.
</details>

#### Additional findings about the router
* `/spec` and `/spec/` are the same, even when the `RedirectTrailingSlash` setting is set to `false`. There are no other configuration options where this behaviour could be tweaked.

### Details of criteria

#### Context type

Custom thing, has everything we need really.

#### Standard library handling

Yep, fasthttp provides an adaptor in the module which we can use to wrap both the handler interface implementations, and the handler funcs as well.

#### Accessing raw Request and ResponseWriter

This is halfway between a kinda and an eeeehhh... It can be done, but not super straightforward. If you look into the function definition of the `fasthttpadaptor.NewFastHTTPHandler` you will see that there's a separate exported function to convert a `fasthttp.Context` into an `http.Request` called `fasthttpadaptor.ConvertRequest`, but assembling the writer is a custom unexported thing that also happens to implement the `http.ResponseWriter` interface. That can be replicated locally, but it's a _pain_, unnecessary extra work, so gonna give it an eh.

#### Websocket

This is a kinda, because there is a [package available](https://github.com/fasthttp/fastws), but it's a fork, but the go.mod file isn't updated, so we have to require it by the original path, at which point the fork isn't getting used, so whatever difference there is between the two is lost.

Also, it's a lot more awkward to work with than the standard library websocket or the gobwas implementation, but it can be done.

#### Path specificity

This one fails our requirements.

* ✅ has catch-all parameter in the form of `/path/{somename:*}`, which needs to go at the end of the route
* ✅ if there's a static route that overlaps with a catch-all, that one gets handled first, regardless of which declaration comes first in the app
* ❌ does not handle `/spec` and `/spec/` differently. Can't make it handle them differently
* ❌ `/path/{somename:*}` does not match `GET /path/` where `somename` would be empty. It either redirects to `/path`, without the trailing slash, or gives us a 404

#### Path variables

Kinda, because even though the context has the `.UserValues(key string)` method, the return type is an `interface{}`, so it's up to us to deal with type checking it, which is inconvenient.

#### Grouping

Works as expected, though it can't do the more advanced things like "attach these middlewares to this group only", so because of that it's a `Kinda`, and not a `Yes`.

#### Overlaps

This one works as expected.

#### General middleware

This is an `Eeeehhh` because even though creating middlewares itself is easy, using them is not really without convenience methods. There's no `.Use()` method or anything similar which would allow us to add middlewares globally, or to a group of our choice, so things like logging or tracing / requestid middlewares would need to be attached to each individual route manually, or create a convenience wrapper curried function that would do it for us.

For example:
```go
type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

func NewMiddlewareWrapper(mws ...Middleware) Middleware {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		for _, mw := range mws {
			h = mw(h)
		}

		return h
    }
}
```
And then use this with:
```go
func main() {
	wr := NewMiddlewareWrapper(Logger(), RequestID(), CORS(), Whatever())

	r.Get("/path", wr(handlers.PathHandler()))
}
```
This is easy, but tedious.

#### Error handling middleware

Divided up into two parts: panic handler, which can be configured on the router, and general error handler, which cannot.

Panic handler just works, the function receives the context and whatever the recovered interface is, and then we're free to do whatever we want, mostly log a lot and send a generic "whoops something went wrong" response back.

The error handler however needs to be a middleware.

The actual error handler middleware is a copy paste from the chi solution with some changes on signatures, but otherwise it's... the same.

You can see how I wrapped everything in [app.go](app/app.go) for the four error middleware routes.

#### Context up and down

Using the `.UserValues()` and `.SetUserValue()` methods on the context it's actually pretty easy.

#### Unit tests

It's weird, because instead of using the `httptest` package, tests are done using the `fasthttp.RequestCtx` struct. That's where we set all the things we would normally configure in an `httptest.NewRequest` call anyways, and then pass that to the handler.

For the assertions the ctx will have its `Response` property hold all the values after a call.

The good thing is that this way we can test the handler itself, and no need to involve the router.

#### Ecosystem

There are a bunch of things made on top of fasthttp, both the handler, and the router. They are linked from the respective repositories, but none of them are the kind of things we could use in conjunction with the router, for example.

They're more a case of "if you want something more full featured, you should use this other thing rather than this thing".

#### Performance

fasthttp router has a router configurable panic handler, but no error handler. The `/smol-perf` request will have the panic handler around it as a result.

##### /performance

* [perftests/run1.log](perftests/run1.log):  Requests/sec:	24845.7544
* [perftests/run2.log](perftests/run2.log):  Requests/sec:	24789.2198
* [perftests/run3.log](perftests/run3.log):  Requests/sec:	25105.3272
* [perftests/run4.log](perftests/run4.log):  Requests/sec:	24716.5910

##### /smol-perf

* [perftests/smol-run1.log](perftests/smol-run1.log):  Requests/sec:	43687.9218
* [perftests/smol-run2.log](perftests/smol-run2.log):  Requests/sec:	41612.6946
* [perftests/smol-run3.log](perftests/smol-run3.log):  Requests/sec:	39553.0190
* [perftests/smol-run4.log](perftests/smol-run4.log):  Requests/sec:	44093.6048
