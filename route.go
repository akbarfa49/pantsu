package pantsu

import (
	"bytes"

	"github.com/cornelk/hashmap"
)

type (
	Route struct {
		Path    string
		Method  string
		Handler Handler
		Params  [][]string
	}
	//routeMap *hashmap.HashMap
)

func buildRoute(url, method string, handler Handler, middleware ...MiddlewareFunc) Route {
	handler = WithMiddlewares(middleware...)(handler)
	return Route{
		Path:    url,
		Method:  method,
		Handler: handler,
	}
}

func (mux *Pantsu) addRoute(r Route) {
	path := s2b(r.Path)
	method := s2b(r.Method)

	if !bytes.ContainsRune(path, ':') {
		if v, ok := mux.routes.Get(path); !ok {
			hm := hashmap.HashMap{}
			hm.Set(method, r)
			mux.routes.Set(path, &hm)
		} else {
			hm := v.(*hashmap.HashMap)
			hm.Set(method, r)
		}
		return
	}
	idx := findPathIndex(r.Path)
	routeParams := path[:idx]
	if v, ok := mux.routes.Get(routeParams); !ok {
		hm := hashmap.HashMap{}
		hm.Set(r.Method, r)
		mux.routes.Set(routeParams, &hm)
	} else {
		hm := v.(*hashmap.HashMap)
		hm.Set(r.Method, r)
	}

}
