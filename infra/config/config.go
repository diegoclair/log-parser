package config

type Config struct {
	AppName  string
	LogDebug bool
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		AppName:  "log-parser",
		LogDebug: true,
	}
}
