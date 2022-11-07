package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// BuildVersion - эта и нижележащие переменные с префиксом build
// являются переменными задаваемыми при компиляции.
var buildVersion string

// Meta - поля с иформацией которые могут пригодиться приложению строго в информативных целях.
type Meta struct {
	SvcVersion string
}

// MockServer - настройки Mock-сервера.
type MockServer struct {
	Port int `default:"8000"`
}

// AppConf общий конфиг приложения.
type AppConf struct {
	Meta       Meta
	MockServer MockServer
}

// NewAppConf - Возвращает настройки заполеннные из переменных окружения.
func NewAppConf() *AppConf {
	cfg := &AppConf{}
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("NewAppConf: %v", err)
	}

	// замена значений на значения заданные при компиляции
	cfg.Meta.SvcVersion = buildVersion

	return cfg
}
