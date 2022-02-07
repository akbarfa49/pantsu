package pantsu_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/akbarfa49/pantsu"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestRealWorld(t *testing.T) {
	cli := http.DefaultClient
	pant := pantsu.NewPantsu(pantsu.Config{
		GlobalErrorHandler: notFoundHandler,
	})
	runGet1 := `ok`
	pant.Get(`/`, func(ctx *fasthttp.RequestCtx) error {
		return pantsu.String(ctx, runGet1)
	})
	go fasthttp.ListenAndServe(`:8080`, pant.ServeHTTP)

	t.Run(`GET OK`, func(t *testing.T) {

		req, err := http.NewRequest(`GET`, `http://localhost:8080/`, nil)

		assert.Empty(t, err)
		if err != nil {
			return
		}

		res, err := cli.Do(req)

		assert.Empty(t, err)
		if err != nil {
			return
		}
		b, _ := io.ReadAll(res.Body)
		fmt.Println(string(b))
		assert.Equal(t, string(b), runGet1)
	})
	t.Run(`GET NOT FOUND`, func(t *testing.T) {

		req, err := http.NewRequest(`GET`, `http://localhost:8080/weirdo`, nil)
		assert.Empty(t, err)
		if err != nil {
			return
		}
		res, err := cli.Do(req)
		assert.Empty(t, err)
		if err != nil {
			return
		}
		assert.Equal(t, res.StatusCode, 404)
	})
}

func notFoundHandler(ctx *fasthttp.RequestCtx) error {
	return pantsu.NotFound(ctx)
}
