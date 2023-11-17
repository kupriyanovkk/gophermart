package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/handlers"
	"github.com/kupriyanovkk/gophermart/internal/middlewares"
	"github.com/kupriyanovkk/gophermart/internal/store"
)

type App struct {
	Flags config.ConfigFlags
	Store store.Store
	// URLChan chan store.DeletedURLs
}

func (a *App) Start() {
	router := chi.NewRouter()

	router.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.JwtAuth,
	)

	router.Route("/api", func(router chi.Router) {
		router.Route("/user", func(router chi.Router) {
			router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
				handlers.Register(w, r, a.Store)
			})

			router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				handlers.Login(w, r, a.Store)
			})

			router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {

			})

			router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
				handlers.PostOrders(w, r, a.Store)
			})

			router.Get("/balance", func(w http.ResponseWriter, r *http.Request) {

			})

			router.Get("/withdrawals", func(w http.ResponseWriter, r *http.Request) {

			})

			router.Post("/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	err := http.ListenAndServe(a.Flags.RunAddress, router)
	if err != nil {
		panic(err)
	}
}

func NewApp(s store.Store, f config.ConfigFlags) *App {
	return &App{Store: s, Flags: f}
}
