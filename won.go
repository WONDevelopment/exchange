package exchange

import (
	"github.com/xiangxian/exchange/pkg"
	"time"
)

type Won interface {
	Time() (time.Time, error)
	Depth(pkg.DepthRequest) (*pkg.DepthResult, error)
	RecentTrades(pkg.TradeRequest) ([]*pkg.RecentTrade, error)
	MyTrades(pkg.TradeRequest) ([]*pkg.MyTrade, error)
	Account(pkg.AccountRequest) (*pkg.Account, error)
	TickerPrice(pkg.TickerPriceRequest) (*pkg.TickerPrice, error)
	CreateOrder(pkg.CreateOrderRequest) (*pkg.Order, error)
	GetOrders(pkg.OrdersRequest) ([]*pkg.Order, error)
	GetOrder(pkg.OrderRequest) (*pkg.Order, error)
	CancelOrder(pkg.CancelOrderRequest) error
}

type won struct {
	Service pkg.Service
}

func NewWon(service pkg.Service) Won {
	return &won{
		Service: service,
	}
}

func (w *won) Time() (time.Time, error) {
	return w.Service.Time()
}
func (w *won) Depth(dr pkg.DepthRequest) (*pkg.DepthResult, error) {
	return w.Service.Depth(dr)
}
func (w *won) RecentTrades(tr pkg.TradeRequest) ([]*pkg.RecentTrade, error) {
	return w.Service.RecentTrades(tr)
}
func (w *won) MyTrades(tr pkg.TradeRequest) ([]*pkg.MyTrade, error) {
	return w.Service.MyTrades(tr)
}
func (w *won) Account(ar pkg.AccountRequest) (*pkg.Account, error) {
	return w.Service.Account(ar)
}
func (w *won) TickerPrice(tpr pkg.TickerPriceRequest) (*pkg.TickerPrice, error) {
	return w.Service.TickerPrice(tpr)
}
func (w *won) CreateOrder(cor pkg.CreateOrderRequest) (*pkg.Order, error) {
	return w.Service.CreateOrder(cor)
}
func (w *won) GetOrders(osr pkg.OrdersRequest) ([]*pkg.Order, error) {
	return w.Service.GetOrders(osr)
}
func (w *won) GetOrder(or pkg.OrderRequest) (*pkg.Order, error) {
	return w.Service.GetOrder(or)
}
func (w *won) CancelOrder(cor pkg.CancelOrderRequest) error {
	return w.Service.CancelOrder(cor)
}
