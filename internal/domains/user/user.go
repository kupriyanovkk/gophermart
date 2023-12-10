package user

import (
	"github.com/kupriyanovkk/gophermart/internal/domains/user/models"
	memStore "github.com/kupriyanovkk/gophermart/internal/domains/user/store/memory"
	pgStore "github.com/kupriyanovkk/gophermart/internal/domains/user/store/pg"
	"github.com/kupriyanovkk/gophermart/internal/shared"
)

type User struct {
	Store models.UserStore
}

func NewUser(db shared.DatabaseConnection) User {
	var store models.UserStore

	if db == nil {
		store = memStore.NewStore()
	} else {
		store = pgStore.NewStore(db)
	}

	order := User{
		Store: store,
	}

	return order
}
