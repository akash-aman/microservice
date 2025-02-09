package models

type Service struct {
	Name     string   `mapstructure:"name" validate:"required"`
	Version  string   `mapstructure:"version" validate:"required"`
	Protocol []string `mapstructure:"protocol" validate:"required"`
}
