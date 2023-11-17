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

func ParseFlags() ConfigFlags {
	var runAddress string
	var databaseURI string
	var accrual string

	flag.StringVar(&runAddress, "a", "localhost:8080", "address and port to run service")
	flag.StringVar(&databaseURI, "d", "postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable", "the address for DB connection")
	flag.StringVar(&accrual, "r", "", "address of the accrual calculation system")

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		runAddress = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		databaseURI = envDatabaseURI
	}
	if envAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrual != "" {
		accrual = envAccrual
	}

	return ConfigFlags{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrual,
	}
}
