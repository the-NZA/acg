package acg

// Config is a base struct for store settings
type Config struct {
	BindAddr    string `json:"bind_addr"`
	LogLevel    string `json:"log_level"`
	DatabaseURL string `json:"db_url"`
}

// NewConfig creates basic config
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		DatabaseURL: "mongodb://localhost:27107",
		LogLevel:    "debug",
	}
}
