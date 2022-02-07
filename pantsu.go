package pantsu

import (
	"bytes"
	"log"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/valyala/fasthttp"
)

type Pantsu struct {
	parent      *Pantsu
	path        string
	routes      *hashmap.HashMap
	routesParam *hashmap.HashMap
	config      Config
	middlewares []MiddlewareFunc
}

func NewPantsu(conf ...Config) *Pantsu {
	pantsu := Pantsu{
		routes:      &hashmap.HashMap{},
		routesParam: &hashmap.HashMap{},
		config:      buildConfig(conf...),
	}

	return &pantsu
}

func (mux *Pantsu) Get(url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		//mux.addRouteFromGroup(buildRoute(url, http.MethodGet, handler))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, fasthttp.MethodGet, handler, middleware...))
}

func (mux *Pantsu) getRoute(ctx *fasthttp.RequestCtx) Route {
	path := ctx.Path()
	method := ctx.Method()
	if mux.config.RemoveTrailingSlash && len(path) > 1 && bytes.HasSuffix(path, []byte(`/`)) {
		path = path[:len(path)-1]
	}
	mux.routes.String()
	foundRouteMap, ok := mux.routes.Get(path)
	if !ok {
		foundRouteMap, _ = mux.routesParam.Get(bFindPathIndex(path))
	}
	var foundRoute Route
	if foundRouteMap == nil {
		foundRoute.Handler = globalErrorHandler
		Error(ctx, 404, `not found`)
		return foundRoute
	}
	route, ok := foundRouteMap.(*hashmap.HashMap)
	if ok {
		v, ok2 := route.Get(method)
		if !ok2 {
			foundRoute.Handler = globalErrorHandler
			Error(ctx, 404, `not found`)
			return foundRoute
		}
		foundRoute = v.(Route)
	}
	return foundRoute
}

func (mux *Pantsu) ServeHTTP(ctx *fasthttp.RequestCtx) {
	tm := time.Now()

	route := mux.getRoute(ctx)
	err := route.Handler(ctx)
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.Write([]byte(err.Error()))
	}
	log.Println(time.Since(tm).Microseconds(), ` microsecond elapsed`)
}

var globalErrorHandler = func(ctx *fasthttp.RequestCtx) error {
	return nil
}
