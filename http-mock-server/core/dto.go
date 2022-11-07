package core

import (
	"net/http"
	"net/url"
)

type RawFixture struct {
	Request struct {
		Method  string              `yaml:"method"`
		Path    string              `yaml:"path"`
		Query   map[string][]string `yaml:"query"`
		Headers map[string][]string `yaml:"headers"`
		Body    string              `yaml:"body"`
	} `yaml:"request"`
	Response struct {
		Code    int                 `yaml:"code"`
		Headers map[string][]string `yaml:"headers"`
		Body    string              `yaml:"body"`
	} `yaml:"response"`
}

type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Headers http.Header
	Body    string
}
type Response struct {
	Code    int
	Headers http.Header
	Body    string
}
type Fixture struct {
	request  Request
	response Response
}
