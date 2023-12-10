package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/gophermart/internal/accrual"
	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance"
	balanceHandlers "github.com/kupriyanovkk/gophermart/internal/domains/balance/handlers"
	"github.com/kupriyanovkk/gophermart/internal/domains/order"
	orderHandlers "github.com/kupriyanovkk/gophermart/internal/domains/order/handlers"
	"github.com/kupriyanovkk/gophermart/internal/domains/user"
	userHandlers "github.com/kupriyanovkk/gophermart/internal/domains/user/handlers"

	"github.com/kupriyanovkk/gophermart/internal/middlewares"
)

func Prepare(flags config.ConfigFlags) {
	db, err := sql.Open("postgres", flags.DatabaseURI)

	if err != nil {
		panic(err)
	}

	accrualChan := make(chan accrual.Accrual)
	accrualClient := accrual.NewClient(flags.AccrualSystemAddress)
	balance := balance.NewBalance(db, accrualChan)
	order := order.NewOrder(db, accrualChan, accrualClient)
	user := user.NewUser(db)

	balanceHandlers.Init(balance.Store)
	orderHandlers.Init(order.Store)
	userHandlers.Init(user.Store)
}

func Start(flags config.ConfigFlags) {
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
				balanceHandlers.GetWithdraw(w, r)
			})

			router.Post("/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
				balanceHandlers.PostWithdraw(w, r)
			})
		})
	})

	err := http.ListenAndServe(flags.RunAddress, router)
	if err != nil {
		panic(err)
	}
}
