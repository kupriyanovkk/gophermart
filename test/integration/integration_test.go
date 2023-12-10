package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
	"github.com/kupriyanovkk/gophermart/internal/app"
	"github.com/kupriyanovkk/gophermart/internal/config"

	"github.com/kupriyanovkk/gophermart/internal/domains/balance"
	balanceHandlers "github.com/kupriyanovkk/gophermart/internal/domains/balance/handlers"
	"github.com/kupriyanovkk/gophermart/internal/domains/order"
	orderHandlers "github.com/kupriyanovkk/gophermart/internal/domains/order/handlers"
	"github.com/kupriyanovkk/gophermart/internal/domains/user"
	userHandlers "github.com/kupriyanovkk/gophermart/internal/domains/user/handlers"

	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
	"github.com/stretchr/testify/suite"
)

type clientTest struct {
}

func (c clientTest) CheckStatus(orderID int) (accrual.Accrual, error) {
	return accrual.Accrual{
		UserID:  TestUsers[len(TestUsers)-1],
		Order:   strconv.Itoa(orderID),
		Status:  accrual.StatusProcessed,
		Accrual: float32(105.4),
	}, nil
}

func NewTestClient() accrual.Client {
	return clientTest{}
}

type TestSuite struct {
	suite.Suite
	appContext    context.Context
	done          context.CancelFunc
	balance       *balance.Balance
	order         *order.Order
	user          *user.User
	accrualClient accrual.Client
	testUsers     []int
	flags         config.ConfigFlags
}

var TestUsers []int

func (suite *TestSuite) initFlags() {
	suite.flags = config.ConfigFlags{
		RunAddress:           "localhost:55000",
		AccrualSystemAddress: "localhost:55001",
	}
}

func (suite *TestSuite) initStorages() {
	accrualChan := make(chan accrual.Accrual)

	balance := balance.NewBalance(nil, accrualChan)
	order := order.NewOrder(nil, accrualChan, suite.accrualClient)
	user := user.NewUser(nil)

	balanceHandlers.Init(balance.Store)
	orderHandlers.Init(order.Store)
	userHandlers.Init(user.Store)

	suite.balance = &balance
	suite.order = &order
	suite.user = &user
}

func (suite *TestSuite) registerTestUsers() {
	type testUser struct {
		Login    string
		Password string
	}

	users := []testUser{
		{Login: "login", Password: "password"},
		{Login: "login1", Password: "password1"},
	}

	for _, user := range users {
		userID, err := suite.user.Store.RegisterUser(suite.appContext, strings.TrimSpace(user.Login), strings.TrimSpace(user.Password))

		suite.Require().NoError(err, "user registration")

		TestUsers = append(TestUsers, userID)
	}
	suite.testUsers = TestUsers
}

var appAlreadyStart = false

func (suite *TestSuite) startApp() {
	if !appAlreadyStart {
		appAlreadyStart = true
		go func() {
			app.Start(suite.flags)
		}()
	}
}

func (suite *TestSuite) initAccrualClient() {
	suite.accrualClient = NewTestClient()
}

func (suite *TestSuite) SetupTest() {
	suite.appContext, suite.done = context.WithCancel(context.Background())
	suite.initFlags()
	suite.initAccrualClient()
	suite.initStorages()
	suite.registerTestUsers()
	suite.startApp()

	time.Sleep(100 * time.Millisecond)
}

func (suite *TestSuite) TearDownTest() {
	suite.done()
}

func (suite *TestSuite) TestUserRegister() {
	client := http.Client{}
	body := strings.NewReader(`{ "login":"login3", "password":"password3" }`)
	url := fmt.Sprintf("http://%s%s", suite.flags.RunAddress, "/api/user/register")
	req, err := http.NewRequest(http.MethodPost, url, body)

	suite.Require().NoError(err, "request error")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	suite.Require().NoError(err, "response error")

	defer res.Body.Close()

	suite.Require().Exactlyf(http.StatusOK, res.StatusCode, "unexpected status")
	authHeader := res.Header.Get("Authorization")
	userID := tokenutil.GetUserIDFromAuthHeader(authHeader)

	_, err = suite.user.Store.GetUser(suite.appContext, userID)
	suite.Require().NoError(err, "user doesn't exist")
}

func (suite *TestSuite) TestUserLogin() {
	client := http.Client{}
	userID := TestUsers[len(TestUsers)-1]
	body := strings.NewReader(`{ "login":"login1", "password":"password1" }`)
	url := fmt.Sprintf("http://%s%s", suite.flags.RunAddress, "/api/user/login")
	req, err := http.NewRequest(http.MethodPost, url, body)

	suite.Require().NoError(err, "request error")

	token := tokenutil.GetBearerHeader(userID)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	suite.Require().NoError(err, "response error")

	defer res.Body.Close()

	suite.Require().Exactlyf(http.StatusOK, res.StatusCode, "unexpected status")
	_, err = suite.user.Store.GetUser(suite.appContext, userID)
	suite.Require().NoError(err, "user doesn't exist")
}

func (suite *TestSuite) TestAddOrder() {
	client := http.Client{}
	orderNumber := `123456789049`
	userID := TestUsers[len(TestUsers)-1]
	body := strings.NewReader(orderNumber)
	url := fmt.Sprintf("http://%s%s", suite.flags.RunAddress, "/api/user/orders")
	req, err := http.NewRequest(http.MethodPost, url, body)

	suite.Require().NoError(err, "request error")

	req.Header.Set("Content-Type", "text/plain")
	token := tokenutil.GetBearerHeader(userID)
	req.Header.Set("Authorization", token)

	res, err := client.Do(req)
	suite.Require().NoError(err, "response error")

	defer res.Body.Close()

	suite.Require().Exactlyf(http.StatusAccepted, res.StatusCode, "unexpected status")
	orders, err := suite.order.Store.GetOrders(suite.appContext, userID)
	suite.Require().NoError(err, "error checking the availability of an order in the db")
	suite.Equal(1, len(orders), "order quantity error")
}

func (suite *TestSuite) TestUserBalance() {
	client := http.Client{}
	orderNumber := `123456789049`
	userID := TestUsers[len(TestUsers)-1]
	body := strings.NewReader(orderNumber)
	url := fmt.Sprintf("http://%s%s", suite.flags.RunAddress, "/api/user/orders")
	req, err := http.NewRequest(http.MethodPost, url, body)

	suite.Require().NoError(err, "request error")

	req.Header.Set("Content-Type", "text/plain")
	token := tokenutil.GetBearerHeader(userID)
	req.Header.Set("Authorization", token)

	res, err := client.Do(req)
	suite.Require().NoError(err, "response error")

	defer res.Body.Close()

	suite.Require().Exactlyf(http.StatusAccepted, res.StatusCode, "unexpected status")
	orders, _ := suite.order.Store.GetOrders(suite.appContext, userID)
	order := orders[0]

	time.Sleep(1 * time.Second)

	userBalance, _ := suite.balance.Store.GetUserBalance(suite.appContext, order.UserID)
	suite.Require().Exactlyf(float32(105.4), userBalance.Current, "incorrect current balance")
}

func (suite *TestSuite) TestBalanceChangedAfterOrderAndWithdraw() {
	suite.TestUserBalance()
	client := http.Client{}
	userID := TestUsers[len(TestUsers)-1]
	bodyString := `{ "order": "2377225624", "sum": 50 }`
	body := strings.NewReader(bodyString)
	url := fmt.Sprintf("http://%s%s", suite.flags.RunAddress, "/api/user/balance/withdraw")
	req, err := http.NewRequest(http.MethodPost, url, body)

	suite.Require().NoError(err, "request error")

	token := tokenutil.GetBearerHeader(userID)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	suite.Require().NoError(err, "response error")

	defer res.Body.Close()

	suite.Require().Exactlyf(http.StatusOK, res.StatusCode, "unexpected status")

	orders, _ := suite.order.Store.GetOrders(suite.appContext, userID)
	order := orders[0]
	userBalance, _ := suite.balance.Store.GetUserBalance(suite.appContext, order.UserID)
	suite.Require().Exactlyf(float32(50), userBalance.Withdrawn, "incorrect withdrawn balance")
}

func TestApp(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
