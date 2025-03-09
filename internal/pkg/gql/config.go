package gql

import (
	"net/http"
	"time"
)

type GraphQLConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	BaseRoute    string        `mapstructure:"baseRoute"`
	DebugMode    bool          `mapstructure:"debugMode"`
}

func NewGQLServer() *http.Server {
	return &http.Server{}
}
