// Package config - единая точка конфигурирования приложения
// Настройки формируют дерево компонентов.
// Все настройки берутся только из переменных окружения
// Настройки по умолчанию не должны приводить к ошибкам сборки
//
// Нельзя передавать весь AppConf в компоненты.
// Стоит ограничиться передачей минимально необходимых подструктур
//
// Структуру настроек сервисов нужно называть {ServiceName}Service
// В структуре Services поле с настройками сервиса называть только по имени сервиса
package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	gb "github.com/sony/gobreaker"
)

// BuildVersion - эта и нижележащие переменные с префиксом build
// являются переменными задаваемые при компиляции.
var buildVersion string

// Meta - поля с иформацией которые могут пригодиться приложению строго в информативных целях.
type Meta struct {
	SvcVersion string
}

type ServiceSettings struct {
	Host           string // scheme://host[:port]
	BasePath       string
	ClientSettings HTTPClientSettings
	UseCB          bool
	Cbs            gb.Settings
}

type HTTPClientSettings struct {
	Timeout           time.Duration `default:"10s"`
	TransportSettings HTTPTransportSettings
}

type HTTPTransportSettings struct {
	DisableKeepAlive bool `default:"true"`
}

// GRPCServer - настройки gRPC сервера.
type GRPCServer struct {
	Port string `default:"50051"`
}

// Metrics - настройки HTTP сервера для метрик prometheus.
type Metrics struct {
	Port             string        `default:"50052"`
	HTTPWriteTimeout time.Duration `default:"10s"`
	HTTPReadTimeout  time.Duration `default:"10s"`
}

// Logger - настройки логирования.
type Logger struct {
	Mode string `default:"prod"`
}

// AppConf общий конфиг приложения.
type AppConf struct {
	Meta       Meta
	GRPCServer GRPCServer
	Metrics    Metrics
	Logger     Logger
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
