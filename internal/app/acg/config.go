package acg

// Config is a base struct for store settings
type Config struct {
	BindAddr    string `json:"bind_addr"`
	LogLevel    string `json:"log_level"`
	DatabaseURL string `json:"db_url"`
	SecretKey   string `json:"secret_key"`
}

// NewConfig creates basic config
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":9999",
		DatabaseURL: "mongodb://test:33333",
		LogLevel:    "debug",
		SecretKey:   "simple_secret",
	}
}
