package main

import (
	"database/sql"

	"github.com/kupriyanovkk/gophermart/internal/app"
	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/store"
)

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
