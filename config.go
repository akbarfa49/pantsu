package pantsu

import "github.com/valyala/fasthttp"

type Config struct {
	RemoveTrailingSlash bool
	GlobalErrorHandler  Handler
}

func buildConfig(conf ...Config) Config {
	config := Config{
		RemoveTrailingSlash: true,
	}
	if len(conf) > 0 {
		config = conf[0]
	}
	if config.GlobalErrorHandler == nil {
		config.GlobalErrorHandler = globalErrorHandler
	} else {

		wrapErrHandler := func(ctx *fasthttp.RequestCtx) error {
			ctx.Response.Reset()
			return config.GlobalErrorHandler(ctx)
		}
		config.GlobalErrorHandler = wrapErrHandler
	}
	return config
}
