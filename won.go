package exchange

import(
	"github.com/xiangxian/exchange/pkg"
	"time"
)

type Won interface {
	Time() (time.Time, error)
	Depth() (pkg.DepthResult, error)
	Trades() (pkg.Trades, error)
	HistoricalTrades()(pkg.Trades, error)
	Account()(pkg.Account, error)
	TickerPrice()(pkg.TickerPrice, error)
	CreateOrder()(pkg.Order, error)
	GetOrders()([]*pkg.Order, error)
	GetOrder()(pkg.Order, error)
	CancelOrder()error
}

type won struct{
	Service pkg.Service
}

func NewWon(service pkg.Service)Won{
	return &won{
		Service: service,
	}
}

func (w *won)Time() (time.Time, error){
	return w.Service.Time()
}
func (w *won)Depth() (pkg.DepthResult, error){
	return w.Service.Depth()
}
func (w *won)Trades() (pkg.Trades, error){
	return w.Service.Trades()
}
func (w *won)HistoricalTrades()(pkg.Trades, error){
	return w.Service.HistoricalTrades()
}
func (w *won)Account()(pkg.Account, error){
	return pkg.Account{}, nil
}
func (w *won)TickerPrice()(pkg.TickerPrice, error){
	return pkg.TickerPrice{}, nil
}
func (w *won)CreateOrder()(pkg.Order, error){
	return w.Service.CreateOrder()
}
func (w *won)GetOrders()([]*pkg.Order, error){
	return w.Service.GetOrders()
}
func (w *won)GetOrder()(pkg.Order, error){
	return w.Service.GetOrder()
}
func (w *won)CancelOrder()error{
	return w.Service.CancelOrder()
}
