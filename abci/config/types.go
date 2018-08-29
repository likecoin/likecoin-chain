package config

// Config is a struct for global configuration
type Config struct {
	Environment string `mapstructure:"env"`
	LogConfig   string `mapstructure:"log_config"`
}

// IsProduction returns true in environment set to production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
