package httpclient

import (
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Headers http.Header
	Body    io.Reader
}
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}
