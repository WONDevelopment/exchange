package pkg

type DepthRequest struct {
	Market string
	Limit  int
}

type TradeRequest struct {
	Market string
	Limit  int
	FromId int64
}

type TickerPriceRequest struct {
	Market string
}

type AccountRequest struct {
	RecvWindow int
	Timestamp  int64
}

type CreateOrderRequest struct {
	Market     string
	Side       string
	Volume     string
	Price      string
	OrdType    string
	RecvWindow int
	Timestamp  int64
}

type OrderRequest struct {
	Id         int64
	RecvWindow int
	Timestamp  int64
}

type CancelOrderRequest struct {
	Id         int64
	RecvWindow int
	Timestamp  int64
}

type OrdersRequest struct {
	Market       string
	OrderId      int64
	State        string
	Side         string
	StartAtStamp int64
	EndAtStamp   int64
	Limit        int
	RecvWindow   int
	Timestamp    int64
}

type DepthResult struct {
	Time int
	Bids []struct {
		Price  string
		Amount string
	}
	Asks []struct {
		Price  string
		Amount string
	}
}

type RecentTrade struct {
	Id       int64
	Price    string
	Quantity string
	CreateAt int64
}

type MyTrade struct {
	Id       int64
	OrderId  int64
	Price    string
	Quantity string
	Side     string
	CreateAt int64
}

type CurrencyAccount struct {
	Currency     string
	TotalBalance string
	Balance      string
	locked       string
	UsdPrice     string
	Precision    int
	Limits       struct {
		MinimalTradeFee string
	}
}

type Account struct {
	Accounts      []CurrencyAccount
	EqualTotalUsd string
}
type TickerPrice struct {
	Market string
	Price  string
}

type Order struct {
	Id              int64  `json:"id"`
	Side            string `json:"side"`
	OrdType         string `json:"ord_type"`
	Price           string `json:"price"`
	State           string `json:"state"`
	Market          string `json:"market"`
	BidCurrency     string `json:"bid_currency"`
	AskCurrency     string `json:"ask_currency"`
	CreatedAtStamp  int64  `json:"created_at_stamp"`
	Volume          string `json:"volume"`
	RemainingVolume string `json:"remaining_volume"`
	ExecutedRate    string `json:"executed_rate"`
	Funds           string `json:"funds"`
}
