package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/gophermart/internal/config"
	balance "github.com/kupriyanovkk/gophermart/internal/domains/balance/handlers"
	order "github.com/kupriyanovkk/gophermart/internal/domains/order/handlers"
	user "github.com/kupriyanovkk/gophermart/internal/domains/user/handlers"
	"github.com/kupriyanovkk/gophermart/internal/middlewares"
	"github.com/kupriyanovkk/gophermart/internal/shared"
)

func init() {
	flags := config.Get()
	db, err := sql.Open("postgres", flags.DatabaseURI)

	if err != nil {
		panic(err)
	}

	loyaltyChan := make(chan shared.LoyaltyOperation)

	balance.Init(db, loyaltyChan)

	order.Init(db, loyaltyChan)

	user.Init(db)
}

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
				user.Register(w, r)
			})

			router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				user.Login(w, r)
			})

			router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
				order.GetOrders(w, r)
			})

			router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
				order.PostOrders(w, r)
			})

			router.Get("/balance", func(w http.ResponseWriter, r *http.Request) {
				balance.GetUserBalance(w, r)
			})

			router.Get("/withdrawals", func(w http.ResponseWriter, r *http.Request) {
				balance.GetWithdraw(w, r)
			})

			router.Post("/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
				balance.PostWithdraw(w, r)
			})
		})
	})

	err := http.ListenAndServe(flags.RunAddress, router)
	if err != nil {
		panic(err)
	}
}
