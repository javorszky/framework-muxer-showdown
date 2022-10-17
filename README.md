# framework-muxer-showdown
A companion repository to the muxer/framework showdown notion page.

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

#### Overlaps

This one works as expected.

#### General middleware

#### Error handling middleware

#### Context up and down

#### Unit tests

#### Ecosystem
