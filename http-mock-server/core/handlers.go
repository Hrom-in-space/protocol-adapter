package core

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"
)

func writeResponse(ctx *fasthttp.RequestCtx, fixtureResponse Response) {
	ctx.Response.SetStatusCode(fixtureResponse.Code)
	for name, values := range fixtureResponse.Headers {
		for _, value := range values {
			ctx.Response.Header.Set(name, value)
		}
	}
	fmt.Fprint(ctx, fixtureResponse.Body)
}

func makeQueryValues(ctx *fasthttp.RequestCtx) url.Values {
	query := url.Values{}
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		query.Add(string(key), string(value))
	})

	return query
}
func makeHeaders(ctx *fasthttp.RequestCtx) http.Header {
	headers := http.Header{}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers.Add(string(key), string(value))
	})

	return headers
}

func NewHandler(mock *Mock) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		query := makeQueryValues(ctx)
		headers := makeHeaders(ctx)
		response, ok := mock.Request(string(ctx.Method()), string(ctx.Path()), query, headers, ctx.Request.Body())
		if ok {
			writeResponse(ctx, response)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotImplemented)
		}
	}
}
