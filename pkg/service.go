package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/xiangxian/exchange"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Service interface {
	Time() (time.Time, error)
	Depth() (exchange.DepthResult, error)
	Trades() (exchange.Trades, error)
	historicalTrades()(exchange.Trades, error)
	Account()(exchange.Account, error)
	TickerPrice()(exchange.TickerPrice, error)
	CreateOrder()(exchange.Order, error)
	GetOrders()([]*exchange.Order, error)
	GetOrder()(exchange.Order, error)
	CancelOrder()error
}

type wonService struct{
	URL string
	APIKey string
	Signer Signer
	Logger log.Logger
	Ctx context.Context
}

func NewWonService(url, apiKey string, signer Signer, logger log.Logger, ctx context.Context) Service {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &wonService{
		URL:    url,
		APIKey: apiKey,
		Signer: signer,
		Logger: logger,
		Ctx:    ctx,
	}
}

func(ws *wonService) Time()(time.Time, error){
	return time.Now(), nil
}

func(ws *wonService) Depth()(exchange.DepthResult, error){
	return exchange.DepthResult{}, nil
}
func(ws *wonService) Trades()(exchange.Trades, error){
	return  exchange.Trades{}, nil
}
func(ws *wonService) historicalTrades()(exchange.Trades, error){
	return  exchange.Trades{}, nil
}
func(ws *wonService) Account()(exchange.Account, error){
	return exchange.Account{}, nil
}
func(ws *wonService) TickerPrice()(exchange.TickerPrice, error){
	return exchange.TickerPrice{}, nil
}
func(ws *wonService) CreateOrder()(exchange.Order, error){
	return exchange.Order{}, nil
}
func(ws *wonService) GetOrders()([]*exchange.Order, error){
	return []*exchange.Order{}, nil
}
func(ws *wonService) GetOrder()(exchange.Order, error){
	return exchange.Order{}, nil
}
func(ws *wonService) CancelOrder()error{
	return nil
}

func (ws *wonService) request(method string, endpoint string, params map[string]string,
	apiKey bool, sign bool) (*http.Response, error) {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("%s/%s", ws.URL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create request error:%s", err.Error()))
	}
	req.WithContext(ws.Ctx)

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	if apiKey {
		req.Header.Add("X-MBX-APIKEY", ws.APIKey)
	}
	if sign {
		level.Debug(ws.Logger).Log("queryString", q.Encode())
		q.Add("signature", ws.Signer.Sign([]byte(q.Encode())))
		level.Debug(ws.Logger).Log("signature", ws.Signer.Sign([]byte(q.Encode())))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
