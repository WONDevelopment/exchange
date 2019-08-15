package exchange

import(
	"github.com/xiangxian/exchange/pkg"
	"time"
)

type Won interface {
	Time() (time.Time, error)
	Depth() (DepthResult, error)
	Trades() (Trades, error)
	historicalTrades()(Trades, error)
	Account()(Account, error)
	TickerPrice()(TickerPrice, error)
	CreateOrder()(Order, error)
	GetOrders()([]*Order, error)
	GetOrder()(Order, error)
	CancelOrder()error
}

type won struct{
	Service pkg.Service
}




type DepthResult struct{}
type Trades struct{}

type currencyAccount struct{}
type Account struct{
	Accounts []currencyAccount
	EqualTotalUsd string
}
type TickerPrice struct{
	Market string
	Price string
}
type Order  struct{}

