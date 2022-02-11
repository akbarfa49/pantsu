package pantsu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestPantsu_Group(t *testing.T) {

	pantsus := NewPantsu()
	pantsu := pantsus.Group(`/1`, Counter())
	pantsu = pantsu.Group(`/2`, Counter())
	pantsu = pantsu.Group(`/3`, Counter())
	pantsu.Get(`/4`, func(ctx *fasthttp.RequestCtx) error {
		ctx.SetStatusCode(ctx.Response.StatusCode() + 1)
		return nil
	})
	ctx := new(fasthttp.RequestCtx)
	ctx.URI().SetPath(`/1/2/3/4`)
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	pantsus.ServeHTTP(ctx)
	assert.Equal(t, 201, ctx.Response.StatusCode())

}

func Counter() MiddlewareFunc {
	return func(next Handler) Handler {
		return func(ctx *fasthttp.RequestCtx) error {
			ctx.SetStatusCode(ctx.Response.StatusCode() + 1)
			return nil
		}
	}
}
