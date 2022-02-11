package pantsu

import (
	"github.com/cornelk/hashmap"
	"github.com/valyala/fasthttp"
)

var (
	_ Handler
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
		mux.addRouteFromGroup(buildRoute(url, fasthttp.MethodGet, handler, middleware...))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, fasthttp.MethodGet, handler, middleware...))
}
func (mux *Pantsu) Post(url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, fasthttp.MethodPost, handler, middleware...))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, fasthttp.MethodPost, handler, middleware...))
}

func (mux *Pantsu) Put(url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, fasthttp.MethodPut, handler, middleware...))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, fasthttp.MethodPut, handler, middleware...))
}

func (mux *Pantsu) Delete(url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, fasthttp.MethodDelete, handler, middleware...))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, fasthttp.MethodDelete, handler, middleware...))
}

func (mux *Pantsu) Patch(url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, fasthttp.MethodPatch, handler, middleware...))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, fasthttp.MethodPatch, handler, middleware...))
}

func (mux *Pantsu) All(url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		for _, method := range methodList {
			mux.addRouteFromGroup(buildRoute(url, method, handler))
		}
		return
	}

	for _, method := range methodList {
		middleware = append(middleware, mux.middlewares...)
		mux.addRoute(buildRoute(url, method, handler, middleware...))
	}
}

//im not sure this will working when using web. but atleast this working when use any http client
func (mux *Pantsu) Custom(method string, url string, handler Handler, middleware ...MiddlewareFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, method, handler, middleware...))
		return
	}
	middleware = append(middleware, mux.middlewares...)
	mux.addRoute(buildRoute(url, method, handler, middleware...))
}
func (mux *Pantsu) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
	mux.config.GlobalErrorHandler = WithMiddlewares(mux.middlewares...)(mux.config.GlobalErrorHandler)
}

func (mux *Pantsu) getRoute(ctx *fasthttp.RequestCtx) Route {
	path := ctx.Path()
	method := ctx.Method()
	if lp := len(path); mux.config.RemoveTrailingSlash && lp > 1 && path[lp-1] == '/' {
		path = path[:lp-1]
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

	route := mux.getRoute(ctx)

	err := route.Handler(ctx)
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.Write([]byte(err.Error()))
	}
}

var globalErrorHandler = func(ctx *fasthttp.RequestCtx) error {
	return nil
}
