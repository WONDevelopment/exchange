package pkg

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