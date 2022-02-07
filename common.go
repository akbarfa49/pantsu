package pantsu

import (
	"github.com/valyala/fasthttp"
)

func WithMiddlewares(middleware ...MiddlewareFunc) MiddlewareFunc {
	return func(next Handler) Handler {
		h := next
		if len(middleware) <= 0 {
			return h
		}

		for i := len(middleware) - 1; i >= 0; i-- {
			if middleware[i] == nil {
				continue
			}
			h = middleware[i](h)
		}
		return h
	}
}

func String(ctx *fasthttp.RequestCtx, str string) (err error) {

	_, err = ctx.WriteString(str)
	return
}

func Error(ctx *fasthttp.RequestCtx, status int, msg string) {
	ctx.SetStatusCode(status)
	ctx.WriteString(msg)

}