//go:generate minimock -g -o ./mmock/ -s .go
package httpclient

import (
	"context"
)

// Executor - интерфейс для выполнения HTTP запросов.
type Executor interface {
	Do(ctx context.Context, req *Request) (*Response, error)
}
