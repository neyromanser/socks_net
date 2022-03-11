package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
	"log"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	ControlDomain    	string `mapstructure:"CONTROL_DOMAIN"`
	AgentPassword    	string `mapstructure:"AGENT_PASSWORD"`
}


// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func GetConfig(path string) (config Config) {

	c, err := LoadConfig(path)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	return c
}

func Dump(some string)  {
	spew.Dump(some)
	return
}