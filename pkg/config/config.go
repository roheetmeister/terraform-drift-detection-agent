package config

// Config holds runtime configuration for a drift scan.
type Config struct {
	StatePath    string
	Provider     string
	Region       string
	OutputFormat string // "table" or "json"
	Port         int
	CronExpr     string
	MaxReports   int
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Provider:     "aws",
		Region:       "us-east-1",
		OutputFormat: "table",
		Port:         8080,
		MaxReports:   50,
	}
}
