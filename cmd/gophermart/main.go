package main

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kupriyanovkk/gophermart/internal/app"
	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/store"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
}

func main() {
	flags := config.ParseFlags()
	db, err := sql.Open("postgres", flags.DatabaseURI)

	if err != nil {
		panic(err)
	}

	store := store.NewStore(db)
	app := app.NewApp(*store, flags)

	app.Start()
}
