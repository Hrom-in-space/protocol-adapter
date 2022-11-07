package core

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	testFixture := `
request:
  # Первая строка HTTP запроса формируется из: {method} {url} {http}
  method: "GET"  # Обязательное
  path: "/path/to/"  # Обязательное
  # Данные добавляемые к url. Имеют приоритет над теми которые указаны в url
  query:
    asd: ["asd"]

  # Заголовки которые имеют приоритет над всеми другими автовычислимыми
  headers:
    key1: ["val1", "val2"]
    key2: ["val3"]

  # Чистые данные
  body: "data"

response:
  code: 201
  # Заголовки которые имеют приоритет над всеми другими автовычислимыми
  headers:
    key1: ["val1", "val2"]
    key2: ["val3"]
  # Чистые данные
  body: "not json"
`
	tests := []struct {
		name string
		want Fixture
		err  error
	}{
		{
			"loaded",
			Fixture{
				request: Request{
					"GET",
					"/path/to/",
					url.Values{"asd": []string{"asd"}},
					http.Header{"key1": []string{"val1", "val2"}, "key2": []string{"val3"}},
					"data",
				},
				response: Response{
					Code:    201,
					Headers: http.Header{"key1": []string{"val1", "val2"}, "key2": []string{"val3"}},
					Body:    "not json",
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := load(strings.NewReader(testFixture))
			require.Equal(t, nil, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
