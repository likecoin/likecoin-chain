package config

// Config is a struct for global configuration
type Config struct {
	Environment    string `mapstructure:"env"`
	LogConfig      string `mapstructure:"log_config"`
	InitialBalance string `mapstructure:"initial_balance"`
	KeepBlocks     int64  `mapstructure:"keep_blocks"`
	DBCacheSize    int    `mapstructure:"db_cache_size"`
}

// IsProduction returns true in environment set to production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
