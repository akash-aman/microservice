package gql

import "time"

type GraphQLConfig struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	BaseRoute    string        `mapstructure:"baseRoute"`
	DebugMode    bool          `mapstructure:"debugMode"`
}
