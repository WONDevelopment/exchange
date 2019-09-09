package tests

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/xiangxian/exchange"
	"github.com/xiangxian/exchange/pkg"
	"net/http"
	"testing"
	"time"
)

func initWon() exchange.Won {
	signer := &pkg.HmacSigner{Key: []byte("your secret key")}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	return exchange.NewWon(server)
}

func TestTime(t *testing.T) {
	won := initWon()
	m, err := won.Time()
	t.Logf(m.String())
	assert.Equal(t, nil, err)
	assert.Equal(t, true, m.UnixNano() < time.Now().UnixNano())
}

func TestDepth(t *testing.T) {
	won := initWon()
	r, err := won.Depth(pkg.DepthRequest{Market: "wonbtc", Limit: 10})
	t.Logf(fmt.Sprintf("DepthResult:%v", r))
	assert.Equal(t, nil, err)
}

func TestTrades(t *testing.T) {
	won := initWon()
	r, err := won.RecentTrades(pkg.TradeRequest{Market: "wonbtc", Limit: 10})
	t.Logf(fmt.Sprintf("TradeResult:%v", r))
	assert.Equal(t, nil, err)
}

func TestHistoryTrades(t *testing.T) {
	won := initWon()
	r, err := won.MyTrades(pkg.TradeRequest{Market: "wonbtc", Limit: 10})
	t.Logf(fmt.Sprintf("HistoricalTrades:%v", r))
	assert.Equal(t, nil, err)
}

func TestTickerPrice(t *testing.T) {
	won := initWon()
	r, err := won.TickerPrice(pkg.TickerPriceRequest{Market: "wonbtc"})
	t.Logf(fmt.Sprintf("TradeResult:%v", r))
	assert.Equal(t, nil, err)
}

func TestAccount(t *testing.T) {
	won := initWon()
	r, err := won.Account(pkg.AccountRequest{Timestamp: int64(time.Now().Unix() * 100), RecvWindow: 5000})
	t.Logf(fmt.Sprintf("AccountResult:%v", r))
	assert.Equal(t, nil, err)
}

func TestGetOrder(t *testing.T) {
	won := initWon()
	r, err := won.GetOrder(pkg.OrderRequest{Id: 1495, Timestamp: int64(time.Now().Unix() * 100), RecvWindow: 5000})
	t.Logf(fmt.Sprintf("GetOrder:%+v", r))
	assert.Equal(t, nil, err)
}

func TestCreateOrder(t *testing.T) {
	won := initWon()
	r, err := won.CreateOrder(pkg.CreateOrderRequest{
		Market:     "wonbtc",
		Side:       "buy",
		Price:      "0.00011",
		Volume:     "1",
		OrdType:    "limit",
		Timestamp:  int64(time.Now().Unix() * 100),
		RecvWindow: 5000})
	t.Logf(fmt.Sprintf("CreateOrder:%+v", r))
	assert.Equal(t, nil, err)
}

func TestCancelOrder(t *testing.T) {
	won := initWon()
	err := won.CancelOrder(pkg.CancelOrderRequest{Id: 1501, Timestamp: int64(time.Now().Unix() * 100), RecvWindow: 5000})

	assert.Equal(t, nil, err)
}

func TestGetOrders(t *testing.T) {
	won := initWon()
	orders, err := won.GetOrders(pkg.OrdersRequest{
		Market:     "wonbtc",
		State:      "wait",
		Side:       "sell",
		RecvWindow: 5000,
		Timestamp:  11110,
	})
	for _, v := range orders {
		t.Logf(fmt.Sprintf("GetOrders:%+v", *v))
	}

	assert.Equal(t, nil, err)
}


func TestSignature(t *testing.T){
	signer := &pkg.HmacSigner{Key: []byte("3a8c1faf-c5c4-4d4d-bb4b-bd42b9e86772")}

	req, _ := http.NewRequest("POST", "api", nil)
	q := req.URL.Query()
	params:= map[string]string{
		"a":"111",
		"b":"111",
	}
	for key, val := range params {
		q.Add(key, val)
	}
	keys:=q.Encode()
	t.Log(keys)
	assert.Equal(t,
		"5877cccfd2dc1328e34c29238d1b54ebaf0d1038ee98d6bd3588f4569fcd42dd",

		signer.Sign([]byte(keys)))
}