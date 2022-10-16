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
</details>

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
