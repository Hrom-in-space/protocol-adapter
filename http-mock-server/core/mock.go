package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func NewMock(fixtures []Fixture) *Mock {
	return &Mock{
		router: makeRouter(fixtures),
	}
}

var nonHashableHeaders = map[string]bool{
	"Connection":      true,
	"User-Agent":      true,
	"Accept-Encoding": true,
	"Host":            true,
	"Accept":          true,
	"Content-Length":  true,
}

type Mock struct {
	router map[string]Response
}

func (mock *Mock) Request(method, path string, query url.Values, headers http.Header, body []byte) (Response, bool) {
	for key := range headers {
		if _, ok := nonHashableHeaders[key]; ok {
			headers.Del(key)
		}
	}
	hash := requestHash(method, path, query, headers, body)
	response, ok := mock.router[hash]

	return response, ok
}

func makeRouter(fixtures []Fixture) map[string]Response {
	data := make(map[string]Response, len(fixtures))
	for _, fixture := range fixtures {
		hash := requestHash(fixture.request.Method, fixture.request.Path, fixture.request.Query, fixture.request.Headers, []byte(fixture.request.Body))
		data[hash] = fixture.response
	}

	return data
}

func headersHash(headers http.Header) []byte {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	buf := bytes.NewBuffer(nil)
	for _, name := range keys {
		buf.WriteString(strings.ToLower(name))
		vals := headers[name]
		sort.Strings(vals)
		for _, v := range vals {
			buf.WriteString(v)
		}
	}

	return buf.Bytes()
}

func requestHash(method, path string, query url.Values, headers http.Header, body []byte) string {
	hash := sha256.New()
	hash.Write([]byte(method))
	hash.Write([]byte(path))
	hash.Write([]byte(query.Encode()))
	hash.Write(headersHash(headers))
	hash.Write(body)

	return hex.EncodeToString(hash.Sum(nil))
}
