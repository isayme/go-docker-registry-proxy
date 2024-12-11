package src

import (
	"sync"

	"github.com/isayme/go-config"
	"github.com/isayme/go-logger"
)

type Config struct {
	Server ServerConfig `json:"server" yaml:"server"`

	Logger LoggerConfig `json:"logger" yaml:"logger"`

	Routes []RouteConfig `json:"routes" yaml:"routes"`
}

type ServerConfig struct {
	Addr string `json:"addr" yaml:"addr"`
}

type LoggerConfig struct {
	Level  string           `json:"level" yaml:"level"`
	Format logger.LogFormat `json:"format" yaml:"format"`
}

type RouteConfig struct {
	Host     string `json:"host" yaml:"host"`
	Upstream string `json:"upstream" yaml:"upstream"`
}

var once sync.Once
var globalConfig = Config{}

func GetConfig() *Config {
	config.Parse(&globalConfig)
	once.Do(func() {
		globalConfig.Default()

		logger.SetLevel(globalConfig.Logger.Level)
		logger.SetFormat(globalConfig.Logger.Format)

		logger.Debugf("log with level: %s, format %s", globalConfig.Logger.Level, globalConfig.Logger.Format)
	})
	return &globalConfig
}

func (c *Config) Default() {
	if c.Server.Addr == "" {
		c.Server.Addr = ":3000"
	}
}
