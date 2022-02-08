package pantsu

import (
	"encoding/json"

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

func Error(ctx *fasthttp.RequestCtx, status int, msg string) error {
	ctx.SetStatusCode(status)
	ctx.WriteString(msg)
	return nil
}

func NotFound(ctx *fasthttp.RequestCtx) error {
	ctx.Response.Reset()

	return JSON(ctx, 404, map[string]interface{}{
		`statusCode`: 404,
		`error`:      `not found`,
	})
}

func Corrupt(ctx *fasthttp.RequestCtx) error {
	ctx.Response.Reset()
	return JSON(ctx, 500, map[string]interface{}{
		`statusCode`: 500,
		`error`:      `internal server error`,
	})
}

func JSON(ctx *fasthttp.RequestCtx, statusCode int, v interface{}) (err error) {
	ctx.SetContentType(`application/json`)
	ctx.SetStatusCode(statusCode)
	byt, err := json.Marshal(v)
	if err != nil {
		return
	}
	_, err = ctx.Write(byt)

	return
}
