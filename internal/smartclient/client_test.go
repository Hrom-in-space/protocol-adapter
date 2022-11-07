package smartclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	gb "github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"

	appconf "protocol-adapter/internal/config"
	"protocol-adapter/internal/httpclient"
	"protocol-adapter/internal/smartclient/mmock"
)

var testHost = "https://testHost:9040"
var testServiceConf = appconf.ServiceSettings{
	Host: testHost,
	ClientSettings: appconf.HTTPClientSettings{
		Timeout: 10 * time.Second,
	},
	UseCB: true,
}

type Settings struct {
	SvcVersion string `envconfig:"SVC_VERSION" default:"dev"`
}

func TestClient_makeURL(t *testing.T) {
	type fields struct {
		host       string
		httpClient *http.Client
		cb         *gb.CircuitBreaker
	}
	type args struct {
		params *httpclient.Request
	}

	fieldsSet := fields{
		testHost,
		&http.Client{},
		gb.NewCircuitBreaker(gb.Settings{}),
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "without query",
			fields: fieldsSet,
			args: args{
				params: &httpclient.Request{
					Method: "GET",
					Path:   "/path/",
				},
			},
			want: testHost + "/path/",
		},
		{
			name:   "with different keys query",
			fields: fieldsSet,
			args: args{
				params: &httpclient.Request{
					Method: "GET",
					Path:   "/path/",
					Query:  url.Values{"key1": []string{"val1"}, "key2": []string{"val2"}},
				},
			},
			want: testHost + "/path/?key1=val1&key2=val2",
		},
		{
			name:   "with two equal keys query",
			fields: fieldsSet,
			args: args{
				params: &httpclient.Request{
					Method: "GET",
					Path:   "/path/",
					Query:  url.Values{"key1": []string{"val1", "val2"}},
				},
			},
			want: testHost + "/path/?key1=val1&key1=val2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, makeURL(tt.fields.host, tt.args.params.Path, tt.args.params.Query), tt.want)
		})
	}
}

func TestClient_Exec_Request(t *testing.T) {
	mc := minimock.NewController(t)
	CBExecutorMock := mmock.NewCBExecutorMock(mc).
		ExecuteMock.Set(func(req func() (interface{}, error)) (p1 interface{}, err error) {
		return req()
	})

	iHTTPClientMock := mmock.NewIHttpClientMock(mc).DoMock.Inspect(func(req *http.Request) {
		assert.Equal(mc, req.Method, "GET")
		assert.Equal(mc, req.URL.String(), testHost+"/user/?query1=1&query2=max")
		assert.Equal(mc, req.Header, http.Header{
			"User-Agent": []string{"protocol-adapter/tests (unknown user-agent)"},
			"header1":    []string{"zip"},
			"header2":    []string{"default"},
		})
		assert.Equal(mc, req.Body, io.NopCloser(strings.NewReader("test_secret")))
	}).Return(&http.Response{Body: io.NopCloser(strings.NewReader(""))}, nil)

	client := NewClient(testServiceConf, appconf.Meta{SvcVersion: "tests"})
	client.httpClient = iHTTPClientMock
	client.cb = CBExecutorMock
	req := &httpclient.Request{
		Method:  "GET",
		Path:    "/user/",
		Query:   url.Values{"query1": []string{"1"}, "query2": []string{"max"}},
		Headers: http.Header{"header1": []string{"zip"}, "header2": []string{"default"}},
		Body:    strings.NewReader("test_secret"),
	}
	if _, err := client.Do(context.Background(), req); err != nil {
		t.Error("Unexpected error:", err)
	}
}

func TestClient_Exec_Response(t *testing.T) {
	mc := minimock.NewController(t)
	CBExecutorMock := mmock.NewCBExecutorMock(mc).
		ExecuteMock.Set(func(req func() (interface{}, error)) (p1 interface{}, err error) {
		return req()
	})

	iHTTPClientMock := mmock.NewIHttpClientMock(mc).DoMock.Set(
		func(req *http.Request) (rp1 *http.Response, err error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader("test_string_body")),
				StatusCode: 200,
				Header:     http.Header{"key": []string{"val"}},
			}, nil
		})

	client := NewClient(testServiceConf, appconf.Meta{})
	client.httpClient = iHTTPClientMock
	client.cb = CBExecutorMock
	req := &httpclient.Request{
		Method: "GET",
		Path:   "/user/",
		Body:   strings.NewReader(""),
	}
	resp, _ := client.Do(context.Background(), req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Body, []byte("test_string_body"))
	assert.Equal(t, resp.Headers, http.Header{"key": []string{"val"}})
}

func Test_getCBExecutor_UseCBfalse(t *testing.T) {
	got := getCBExecutor(appconf.ServiceSettings{UseCB: false})
	assert.Equal(t, got, &zeroCircuitBreaker{})
}

func Test_getCBExecutor_UseCBtrue(t *testing.T) {
	gbs := gb.Settings{Name: "test"}
	got := getCBExecutor(appconf.ServiceSettings{UseCB: true, Cbs: gbs})
	result := func(cbe CBExecutor) bool {
		switch cbe.(type) {
		case *gb.CircuitBreaker:
			return true
		default:
			return false
		}
	}(got)
	assert.Equal(t, true, result)
}
