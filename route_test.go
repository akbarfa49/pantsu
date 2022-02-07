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

	// mux1 := NewNgamux()
	// mux1.Get("/:a", func(rw http.ResponseWriter, r *http.Request) error {
	// 	return String(rw, "ok")
	// })

	// req1 := httptest.NewRequest(http.MethodGet, "/123", nil)
	// rec1 := httptest.NewRecorder()
	// handler1, req1 := mux1.getRoute(req1)
	// handler1.Handler(rec1, req1)

	// result = strings.ReplaceAll(rec1.Body.String(), "\n", "")
	// expected = "ok"
	// must.Equal(expected, result)

	// req2 := httptest.NewRequest(http.MethodGet, "/123", nil)
	// rec2 := httptest.NewRecorder()
	// handler2, req2 := mux.getRoute(req2)
	// handler2.Handler(rec2, req2)

	// result = strings.ReplaceAll(rec2.Body.String(), "\n", "")
	// expected = "not found"
	// must.Equal(expected, result)

	// req3 := httptest.NewRequest(http.MethodPost, "/", nil)
	// rec3 := httptest.NewRecorder()
	// handler2, req3 = mux.getRoute(req3)
	// handler2.Handler(rec3, req3)

	// result = strings.ReplaceAll(rec3.Body.String(), "\n", "")
	// expected = "method not allowed"
	// must.Equal(expected, result)
}
