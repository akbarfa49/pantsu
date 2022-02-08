package pantsu_test

import (
	"io"
	"log"
	"net/http"
	"testing"
	"time"

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
	pant.Use(timer())
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

func timer() func(next pantsu.Handler) pantsu.Handler {
	return func(next pantsu.Handler) pantsu.Handler {

		return func(ctx *fasthttp.RequestCtx) error {
			t := time.Now()
			next(ctx)
			log.Println(`time until done `, time.Since(t).Nanoseconds(), ` nano second`)
			return nil
		}
	}
}
