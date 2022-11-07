//go:generate minimock -g -o ./mmock/ -s .go
// Package smartclient - реализация клиента httpclient.Executor
package smartclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	gb "github.com/sony/gobreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	appconf "protocol-adapter/internal/config"
	"protocol-adapter/internal/httpclient"
	"protocol-adapter/internal/logger"
)

// IHttpClient - interface for http clients like http.Client.
type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CBExecutor - interface for circuit breaker like CircuitBreaker.
type CBExecutor interface {
	Execute(req func() (interface{}, error)) (interface{}, error)
}

type client struct {
	config     appconf.ServiceSettings
	httpClient IHttpClient
	cb         CBExecutor
	meta       appconf.Meta
}

// makeDefaultClient - создает дефолтный http-клиент
// учитывая параметры в переданной конфигурации.
func makeDefaultClient(config appconf.HTTPClientSettings) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = config.TransportSettings.DisableKeepAlive

	return &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}
}

// zeroCircuitBreaker - болванка для соответствия
// интерфейсу CBExecutor.
type zeroCircuitBreaker struct{}

// данный метод просто вызывает переданную в качестве аргумента функцию
// и возвращает результат её выполнения.
func (zcb *zeroCircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	response, err := req()
	if err != nil {
		return nil, err
	}

	return response.(*httpclient.Response), nil
}

// getCBExecutor - возвращает реальный gobreaker или zeroCircuitBreaker
// в зависимости от настроек клиента.
func getCBExecutor(config appconf.ServiceSettings) CBExecutor {
	var circuitBreaker CBExecutor
	if config.UseCB {
		circuitBreaker = gb.NewCircuitBreaker(config.Cbs)
	} else {
		circuitBreaker = &zeroCircuitBreaker{}
	}

	return circuitBreaker
}

func NewClient(config appconf.ServiceSettings, meta appconf.Meta) *client {
	httpClient := makeDefaultClient(config.ClientSettings)
	cb := getCBExecutor(config)

	return &client{
		config:     config,
		httpClient: httpClient,
		cb:         cb,
		meta:       meta,
	}
}

// TODO: вынести в файл с утилитами.
func makeURL(host string, path string, query url.Values) string {
	fullURL := host + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	return fullURL
}

// uaFromMetadataFromCtx - extract user-agent header form metadata stored in context.
func uaFromMetadataFromCtx(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	ua := "unknown user-agent"

	if ok {
		uaHeaders := md.Get("User-Agent")
		if len(uaHeaders) > 0 {
			ua = uaHeaders[0]
		}
	}

	return ua
}

// makeRequestFunc - create function for CB and closure in they all required params.
func (c *client) makeRequestFunc(ctx context.Context, req *httpclient.Request) func() (interface{}, error) {
	log := ctxzap.Extract(ctx)
	path := c.config.BasePath + req.Path
	requestURL := makeURL(c.config.Host, path, req.Query)
	meta := c.meta
	httpClient := c.httpClient

	return func() (interface{}, error) {
		// create *http.Request
		request, err := http.NewRequest(req.Method, requestURL, req.Body)
		if err != nil {
			return nil, err
		}
		// add Headers
		if len(req.Headers) > 0 {
			request.Header = req.Headers
		}
		request.Header.Add(
			"User-Agent",
			fmt.Sprintf("protocol-adapter/%s (%s)", meta.SvcVersion, uaFromMetadataFromCtx(ctx)),
		)

		// exec request
		resp, err := httpClient.Do(request)
		if err != nil {
			return nil, err
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				logger.Warn(ctx, log, "Close body error:", zap.String("error", err.Error()))
			}
		}()

		// read body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// create Response
		response := &httpclient.Response{
			StatusCode: resp.StatusCode,
			Body:       body,
			Headers:    resp.Header,
		}

		return response, nil
	}
}

// Do - execute http request using setting for service to use or not CircuitBreaker pattern.
func (c *client) Do(ctx context.Context, req *httpclient.Request) (*httpclient.Response, error) {
	response, err := c.cb.Execute(c.makeRequestFunc(ctx, req))
	if err != nil {
		log := ctxzap.Extract(ctx)
		logger.Error(
			ctx, log, "Circuit breaker error", zap.Error(err),
			zap.String("host", c.config.Host),
			zap.String("path", req.Path),
		)

		return nil, err
	}

	return response.(*httpclient.Response), nil
}
