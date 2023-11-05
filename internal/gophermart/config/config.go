package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	ProductionMode       bool   `env:"PRODUCTION_MODE"`
}

func ParseConfig() error {
	flag.StringVar(&Config.RunAddress, "a", ":8080", "service start address and port")
	flag.StringVar(&Config.DatabaseURI, "d", "", "database connection address")
	flag.StringVar(&Config.AccrualSystemAddress, "r", "", "address of the accrual calculation system")
	flag.BoolVar(&Config.ProductionMode, "p", true, "production mode")

	flag.Parse()
	return env.Parse(&Config)
}
