package config

import (
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`

	HmacTokenSecret string `env:"HMAC_TOKEN_SECRET"`
	MigrationsFlag  bool   `env:"MIGRATIONS_FLAG"`
	ProductionMode  bool   `env:"PRODUCTION_MODE"`
}

func ParseConfig() error {
	flag.StringVar(&Config.RunAddress, "a", ":8080", "service start address and port")
	flag.StringVar(&Config.DatabaseURI, "d", "", "database connection address")
	flag.StringVar(&Config.AccrualSystemAddress, "r", "", "address of the accrual calculation system")

	flag.StringVar(&Config.HmacTokenSecret, "h", "developerSecretKey", "hmac for encrypt JWT token")
	flag.BoolVar(&Config.MigrationsFlag, "m", true, "migration flag")
	flag.BoolVar(&Config.ProductionMode, "p", true, "production mode")

	flag.Parse()

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return env.Parse(&Config)
}
