package config

import (
	"github.com/spf13/viper"
)

type ConfigUtility struct {
	// Add fields for configuration settings as needed
	cfg *Config
}

func NewConfigUtility() (*ConfigUtility, error) {

	v := viper.New()

	// reading .env (mandatory)
	v.SetConfigFile(".env")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return &ConfigUtility{cfg: cfg}, nil
}

func (c *ConfigUtility) Config() *Config {
	return c.cfg
}
