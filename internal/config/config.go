package config

import (
	"flag"
	"os"
)

type ConfigFlags struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

var runAddress string
var databaseURI string
var accrualAddress string

func init() {
	parseFlags()
}

func parseFlags() ConfigFlags {
	flag.StringVar(&runAddress, "a", "http://localhost:8081", "address and port to run service")
	flag.StringVar(&databaseURI, "d", "", "the address for DB connection")
	flag.StringVar(&accrualAddress, "r", "http://localhost:8080", "address of the accrual calculation system")

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		runAddress = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		databaseURI = envDatabaseURI
	}
	if envAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrual != "" {
		accrualAddress = envAccrual
	}

	return ConfigFlags{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualAddress,
	}
}

func Get() ConfigFlags {
	return ConfigFlags{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualAddress,
	}
}
