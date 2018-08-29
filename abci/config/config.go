package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	c         *Config
	v         *viper.Viper
	envPrefix = "LIKECHAIN"
)

// GetConfig returns the singleton config
func GetConfig() *Config {
	return c
}

// GetViper returns a viper
func GetViper() *viper.Viper {
	return v
}

func prefixKey(key string) string {
	return strings.ToUpper(envPrefix + "_" + key)
}

func readConfig() {
	v = viper.New()
	v.SetConfigType("toml")

	v.AddConfigPath(".")
	v.AddConfigPath("$GOPATH/src/github.com/likecoin/likechain/abci")
	v.SetConfigName("config")

	// If a config file is found, read it in
	err := v.ReadInConfig()
	if err == nil {
		log.Printf("Using config file:\n\t%s\n", v.ConfigFileUsed())
	} else {
		// Otherwise load default config
		log.Printf("Unable to read config file:\n%v\nUsing default config\n", err)
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
