package config

// Config is a struct for global configuration
type Config struct {
	Environment string `mapstructure:"env"`
	Log         LogConfig
}

// IsProduction returns true in environment set to production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// LogConfig is a struct for logging related configuration
type LogConfig struct {
	Level string
}
