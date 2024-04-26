package config

import "time"

type Config struct {
	Host          string   `env:"HOST,notEmpty"`
	Port          int      `env:"PORT" envDefault:"8080"`
	IsDevelopment bool     `env:"IS_DEVELOPMENT,notEmpty" envDefault:"true"`
	ProxyHeader   string   `env:"PROXY_HEADER" envDefault:"X-Real-IP"`
	LogFields     []string `env:"LOG_FIELDS" envSeparator:","`
	Database      Database
	JwtConfig     JwtConfig
}

type Database struct {
	Driver string `env:"DB_DRIVER" envDefault:"sqlite"`
	DSN    string `env:"DB_DSN" envDefault:"file::memory:?cache=shared"`
}

type JwtConfig struct {
	PrivateKey string        `env:"JWT_PRIVATE_KEY,notEmpty" envDefault:""`
	PublicKey  string        `env:"JWT_PUBLIC_KEY,notEmpty" envDefault:""`
	ExpiresIn  time.Duration `env:"JWT_EXPIRES_IN,notEmpty" envDefault:"24h"`
}
