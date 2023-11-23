package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/gophermart/internal/config"
	balanceHandlers "github.com/kupriyanovkk/gophermart/internal/domains/balance/handlers"
	orderHandlers "github.com/kupriyanovkk/gophermart/internal/domains/order/handlers"
	userHandlers "github.com/kupriyanovkk/gophermart/internal/domains/user/handlers"
	"github.com/kupriyanovkk/gophermart/internal/middlewares"
)

func Start() {
	flags := config.Get()
	router := chi.NewRouter()

	router.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.JwtAuth,
	)

	router.Route("/api", func(router chi.Router) {
		router.Route("/user", func(router chi.Router) {
			router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
				userHandlers.Register(w, r)
			})

			router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				userHandlers.Login(w, r)
			})

			router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
				orderHandlers.GetOrders(w, r)
			})

			router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
				orderHandlers.PostOrders(w, r)
			})

			router.Get("/balance", func(w http.ResponseWriter, r *http.Request) {
				balanceHandlers.GetUserBalance(w, r)
			})

			router.Get("/withdrawals", func(w http.ResponseWriter, r *http.Request) {

			})

			router.Post("/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	err := http.ListenAndServe(flags.RunAddress, router)
	if err != nil {
		panic(err)
	}
}
