package pantsu

import "github.com/valyala/fasthttp"

type (
	MiddlewareFunc func(next Handler) Handler
	Handler        func(ctx *fasthttp.RequestCtx) error
)
