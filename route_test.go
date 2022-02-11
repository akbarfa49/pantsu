package pantsu

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestBuildRoute(t *testing.T) {

	result := buildRoute("/", http.MethodGet, func(ctx *fasthttp.RequestCtx) error { return nil })
	expected := Route{
		Path:   "/",
		Method: http.MethodGet,
	}
	assert.Equal(t, expected.Path, result.Path)
	assert.Equal(t, expected.Method, result.Method)
}

func TestAddRoute(t *testing.T) {
	must := assert.New(t)
	mux := NewPantsu()
	mux.addRoute(buildRoute("/", http.MethodGet, nil))
	mux.addRoute(buildRoute("/a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b", http.MethodGet, nil))
	result := mux.routes.Len()
	expected := 3
	must.Equal(expected, result)

	mux = NewPantsu()
	mux.addRoute(buildRoute("/a/:a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b/:b", http.MethodGet, nil))
	mux.addRoute(buildRoute("/c/:c", http.MethodGet, nil))
	result = mux.routes.Len()
	expected = 3
	must.Equal(expected, result)

}

func TestGetRoute(t *testing.T) {
	must := assert.New(t)
	mux := NewPantsu()

	mux.Get("/", func(ctx *fasthttp.RequestCtx) error {
		return String(ctx, `ok`)
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI(`/`)
	handler := mux.getRoute(ctx)
	handler.Handler(ctx)

	result := string(ctx.Response.Body())
	expected := "ok"
	must.Equal(expected, result)

}

func TestCustomRoute(t *testing.T) {
	must := assert.New(t)
	mux := NewPantsu()

	mux.Custom(`ESTEHPANAS`, "/", func(ctx *fasthttp.RequestCtx) error {
		return String(ctx, `ok`)
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI(`/`)
	ctx.Request.Header.SetMethod(`ESTEHPANAS`)
	handler := mux.getRoute(ctx)
	handler.Handler(ctx)

	result := string(ctx.Response.Body())
	expected := "ok"
	must.Equal(expected, result)

}
