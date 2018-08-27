package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	c         *Config
	envPrefix = "LIKECHAIN"
)

// GetConfig returns the singleton config
func GetConfig() *Config {
	return c
}

func prefixKey(key string) string {
	return strings.ToUpper(envPrefix + "_" + key)
}

func readConfig() {
	v := viper.New()
	v.SetConfigType("toml")

	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.SetConfigName("config")

	// If a config file is found, read it in
	if err := v.ReadInConfig(); err == nil {
		log.Printf("Using config file:\n\t%s\n", v.ConfigFileUsed())
	} else {
		// Otherwise load default config
		log.Println("Using default config")
		v.ReadConfig(bytes.NewBuffer(defaultConf))
	}

	// Read in environment variables that match
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Decode config into struct
	if err := v.Unmarshal(c); err != nil {
		log.Panicf("Unable to decode into struct, %v", err)
	}
}

func init() {
	c = new(Config)

	readConfig()
}
