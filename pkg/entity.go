package pkg

type DepthRequest struct{
	Market string
	Limit int
}

type TradeRequest struct{
	Market string
	Limit int
	FromId int
}

type TickerPriceRequest struct{
	Market string
}

type AccountRequest struct{
	RecvWindow int
	Timestamp int
}

type DepthResult struct{
	Time int
	Bids [] struct{
		Price string
		Amount string
	}
	Asks [] struct{
		Price string
		Amount string
	}
}


type Trade struct{
	Id int
	Price string
	Quantity string
	CreateAt int
}

type HistoryTrade struct{
	Id int
	Price string
	Quantity string
	Side string
	CreateAt int
}

type CurrencyAccount struct{
	Currency string
	TotalBalance string
	Balance string
	locked string
	UsdPrice string
	Precision int
	Limits struct{
		MinimalTradeFee string
	}
}


type Account struct{
	Accounts []CurrencyAccount
	EqualTotalUsd string
}
type TickerPrice struct{
	Market string
	Price string
}
type Order  struct{}