package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/handlers"
	"github.com/kupriyanovkk/gophermart/internal/middlewares"
	"github.com/kupriyanovkk/gophermart/internal/order"
	"github.com/kupriyanovkk/gophermart/internal/store"
	"go.uber.org/zap"
)

type App struct {
	Flags      config.ConfigFlags
	Store      store.Store
	OrdersChan chan int
}

func (a *App) Start() {
	router := chi.NewRouter()

	router.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.JwtAuth,
	)

	go a.checkOrderStatus()

	router.Route("/api", func(router chi.Router) {
		router.Route("/user", func(router chi.Router) {
			router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
				handlers.Register(w, r, a.Store)
			})

			router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				handlers.Login(w, r, a.Store)
			})

			router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetOrders(w, r, a.Store)
			})

			router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
				handlers.PostOrders(w, r, a.Store, a.OrdersChan)
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

func (a *App) checkOrderStatus() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	for orderID := range a.OrdersChan {
		status, err := order.CheckStatus(orderID, a.Flags.AccrualSystemAddress)

		if err != nil {
			sugar.Errorln(
				"err", err.Error(),
				"orderID", orderID,
			)
			return
		}

		if status.Status != order.OrderStatusNotRegister {
			err = a.Store.UpdateOrder(context.TODO(), status)
			if err != nil {
				sugar.Errorln(
					"err", err.Error(),
					"status", status,
				)
				return
			}

			if status.Status != order.OrderStatusProcessed {
				time.AfterFunc(5*time.Second, func() {
					a.OrdersChan <- orderID
				})
			}
		}
	}
}

func NewApp(s store.Store, f config.ConfigFlags) *App {
	return &App{Store: s, Flags: f, OrdersChan: make(chan int, 10)}
}
