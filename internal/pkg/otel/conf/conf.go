package conf

type OtelConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Grpc     string `mapstructure:"grpc"`
	Http     string `mapstructure:"http"`
	Service  string `mapstructure:"service"`
	Insecure bool   `mapstructure:"insecure"`
}
