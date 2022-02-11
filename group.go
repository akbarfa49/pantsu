package pantsu

func (mux *Pantsu) addRouteFromGroup(route Route) {
	middlewares := mux.middlewares
	middlewares = append(middlewares, mux.parent.middlewares...)

	mux.parent.addRoute(buildRoute(mux.path+route.Path, route.Method, route.Handler, middlewares...))
}

func (mux *Pantsu) Group(url string, middlewares ...MiddlewareFunc) *Pantsu {
	if mux.parent == nil {
		group := &Pantsu{
			parent:      mux,
			path:        url,
			middlewares: append(mux.middlewares, middlewares...),
		}
		return group
	}
	group := &Pantsu{
		parent:      mux.parent,
		path:        mux.path + url,
		middlewares: append(mux.middlewares, middlewares...),
	}
	return group
}
