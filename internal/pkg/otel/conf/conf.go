package conf

type OtelConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Grpc     int    `mapstructure:"grpc"`
	Http     int    `mapstructure:"http"`
	Service  string `mapstructure:"service"`
	Insecure bool   `mapstructure:"insecure"`
}
