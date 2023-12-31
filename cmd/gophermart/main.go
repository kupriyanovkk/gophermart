package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/kupriyanovkk/gophermart/internal/app"
	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/shared"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	flags := config.Get()
	db, err := sql.Open("postgres", flags.DatabaseURI)

	if err != nil {
		panic(err)
	}

	shared.BootstrapDB(context.Background(), db)
	app.Prepare(flags)
	app.Start(flags)
}
