package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../config")
	viper.AddConfigPath("$HOME/.checkduplicatedfiles")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return &Config{viper: viper.GetViper()}
}

func (c *Config) Get(key string) string {
	return c.viper.GetString(key)
}

func (c *Config) GetWithDefault(key string, defaultValue string) string {
	val := c.viper.GetString(key)
	if val == "" {
		return defaultValue
	}
	return val
}
func (c *Config) GetListWithDefault(key string, defaultValue string) []string {
	val := c.viper.GetStringSlice(key)
	if len(val) == 0 {
		return strings.Split(defaultValue, ",")
	}
	return val
}
