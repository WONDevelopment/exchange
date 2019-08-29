package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Service interface {
	Time() (time.Time, error)
	Depth(DepthRequest) (*DepthResult, error)
	RecentTrades(TradeRequest) ([]*RecentTrade, error)
	MyTrades(TradeRequest) ([]*MyTrade, error)
	Account(AccountRequest) (*Account, error)
	TickerPrice(TickerPriceRequest) (*TickerPrice, error)
	CreateOrder(CreateOrderRequest) (*Order, error)
	GetOrders(OrdersRequest) ([]*Order, error)
	GetOrder(OrderRequest) (*Order, error)
	CancelOrder(CancelOrderRequest) error
}

type wonService struct {
	URL    string
	APIKey string
	Signer Signer
	Logger log.Logger
	Ctx    context.Context
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

func (ws *wonService) Time() (time.Time, error) {
	params := make(map[string]string)
	res, err := ws.request("GET", "api/v1/time", params, false, false)
	if err != nil {
		return time.Time{}, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("unable to read response from Time:%s", err.Error()))
	}
	defer res.Body.Close()

	type data struct {
		Time int64 `json:time`
	}
	var rawTime struct {
		Date data `json:"data"`
	}
	if err := json.Unmarshal(textRes, &rawTime); err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("Time Response unmarshal failed:%s", err.Error()))
	}
	t, err := timeFromUnixMillTimestamp(rawTime.Date.Time)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (ws *wonService) Depth(dq DepthRequest) (*DepthResult, error) {
	params := make(map[string]string)
	params["market"] = dq.Market
	if dq.Limit > 0 {
		params["limit"] = strconv.Itoa(dq.Limit)
	}

	res, err := ws.request("GET", "api/v1/depth", params, false, false)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from Depth:%s", err.Error()))
	}
	defer res.Body.Close()

	type result struct {
		Time int
		Bids [][]string
		Asks [][]string
	}

	var rawDepth struct {
		Data result `json:"data"`
	}

	if err := json.Unmarshal(textRes, &rawDepth); err != nil {
		return nil, errors.New(fmt.Sprintf("Depth Response unmarshal Depth:%s", err.Error()))
	}
	var resultDepth DepthResult
	type depth struct {
		Price  string
		Amount string
	}
	for _, v := range rawDepth.Data.Bids {
		resultDepth.Bids = append(resultDepth.Bids, depth{Price: v[0], Amount: v[1]})
	}
	for _, v := range rawDepth.Data.Asks {
		resultDepth.Asks = append(resultDepth.Asks, depth{Price: v[0], Amount: v[1]})
	}
	resultDepth.Time = rawDepth.Data.Time

	return &resultDepth, err
}
func (ws *wonService) RecentTrades(tr TradeRequest) ([]*RecentTrade, error) {
	params := make(map[string]string)
	params["market"] = tr.Market
	if tr.Limit > 0 {
		params["limit"] = strconv.Itoa(tr.Limit)
	}

	if tr.FromId > 0 {
		params["from_id"] = strconv.FormatInt(tr.FromId, 10)
	}

	res, err := ws.request("GET", "api/v1/trades/recent", params, true, false)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from Depth:%s", err.Error()))
	}
	defer res.Body.Close()

	type result struct {
		Id       int64  `json:"id"`
		Price    string `json:"price"`
		Quantity string `json:"qty"`
		CreateAt int64  `json:"time"`
	}

	var rawDepth struct {
		Data []result `json:"data"`
	}

	if err := json.Unmarshal(textRes, &rawDepth); err != nil {
		return nil, errors.New(fmt.Sprintf("Trades Response unmarshal Trades:%s", err.Error()))
	}
	var trades []*RecentTrade
	for _, v := range rawDepth.Data {
		trades = append(trades, &RecentTrade{Id: v.Id, Price: v.Price, Quantity: v.Quantity, CreateAt: v.CreateAt})
	}

	return trades, nil
}
func (ws *wonService) MyTrades(tr TradeRequest) ([]*MyTrade, error) {
	params := make(map[string]string)
	params["market"] = tr.Market
	if tr.Limit > 0 {
		params["limit"] = strconv.Itoa(tr.Limit)
	}

	if tr.FromId > 0 {
		params["from_id"] = strconv.FormatInt(tr.FromId, 10)
	}

	res, err := ws.request("GET", "api/v1/trades/my", params, true, true)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from HistoricalTrades:%s", err.Error()))
	}
	defer res.Body.Close()

	type result struct {
		Id       int64  `json:"id"`
		OrderId  int64  `json:"order_id"`
		Price    string `json:"price"`
		Side     string `json:"side"`
		Quantity string `json:"qty"`
		CreateAt int64  `json:"time"`
	}

	var rawDepth struct {
		Data []result `json:"data"`
	}

	if err := json.Unmarshal(textRes, &rawDepth); err != nil {
		return nil, errors.New(fmt.Sprintf("HistoricalTrades Response unmarshal HistoricalTrades:%s", err.Error()))
	}
	var trades []*MyTrade
	for _, v := range rawDepth.Data {
		trades = append(trades, &MyTrade{Id: v.Id, OrderId: v.OrderId, Price: v.Price, Side: v.Side, Quantity: v.Quantity, CreateAt: v.CreateAt})
	}

	return trades, nil
}
func (ws *wonService) Account(ar AccountRequest) (*Account, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(ar.Timestamp, 10)
	if ar.RecvWindow > 0 {
		params["recv_window"] = strconv.Itoa(ar.RecvWindow)
	}
	res, err := ws.request("GET", "api/v1/account", params, true, true)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from Account:%s", err.Error()))
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ws.handleError(textRes)
	}

	type CA struct {
		Currency     string `json:"currency"`
		TotalBalance string `json:"total_balance"`
		Balance      string `json:"balance"`
		locked       string `json:"locked"`
		UsdPrice     string `json:"usd_price"`
		Precision    int    `json:"precision"`
		Limits       struct {
			MinimalTradeFee string `json:"minimal_trade_fee"`
		} `json:"limits"`
	}

	var rawResult struct {
		Data struct {
			Accounts      []CA   `json:"accounts"`
			EqualTotalUsd string `json:"equal_total_usd"`
		} `json:"data"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.New(fmt.Sprintf("Account Response unmarshal failed:%s", err.Error()))
	}

	account := &Account{}
	account.EqualTotalUsd = rawResult.Data.EqualTotalUsd
	for _, v := range rawResult.Data.Accounts {
		account.Accounts = append(account.Accounts, CurrencyAccount{
			Currency:     v.Currency,
			TotalBalance: v.TotalBalance,
			Balance:      v.Balance,
			locked:       v.locked,
			UsdPrice:     v.UsdPrice,
			Precision:    v.Precision,
			Limits:       struct{ MinimalTradeFee string }{MinimalTradeFee: v.Limits.MinimalTradeFee},
		})
	}

	return account, nil
}
func (ws *wonService) TickerPrice(tqr TickerPriceRequest) (*TickerPrice, error) {
	params := make(map[string]string)
	params["market"] = tqr.Market

	res, err := ws.request("GET", "api/v1/ticker/price", params, false, false)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from TickerPrice:%s", err.Error()))
	}
	defer res.Body.Close()

	var rawDepth struct {
		Data TickerPrice `json:"data"`
	}

	if err := json.Unmarshal(textRes, &rawDepth); err != nil {
		return nil, errors.New(fmt.Sprintf("TickerPrice Response unmarshal TickerPrice:%s", err.Error()))
	}

	return &rawDepth.Data, nil
}

func (ws *wonService) CreateOrder(cor CreateOrderRequest) (*Order, error) {
	params := make(map[string]string)
	params["market"] = cor.Market
	params["side"] = cor.Side
	params["volume"] = cor.Volume
	params["price"] = cor.Price
	params["ord_type"] = cor.OrdType
	params["timestamp"] = strconv.FormatInt(cor.Timestamp, 10)
	if cor.RecvWindow > 0 {
		params["recv_window"] = strconv.Itoa(cor.RecvWindow)
	}
	res, err := ws.request("POST", "api/v1/order/create", params, true, true)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from CreateOrder:%s", err.Error()))
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, ws.handleError(textRes)
	}

	var rawResult struct {
		Data Order `json:"data"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.New(fmt.Sprintf("CreateOrder Response unmarshal failed:%s", err.Error()))
	}

	return &rawResult.Data, nil
}
func (ws *wonService) GetOrders(osr OrdersRequest) ([]*Order, error) {
	params := make(map[string]string)
	params["market"] = osr.Market
	params["order_id"] = strconv.FormatInt(osr.OrderId, 10)
	params["start_at_stamp"] = strconv.FormatInt(osr.StartAtStamp, 10)
	params["end_at_stamp"] = strconv.FormatInt(osr.EndAtStamp, 10)
	params["timestamp"] = strconv.FormatInt(osr.Timestamp, 10)

	if len(osr.State) > 0 {
		params["state"] = osr.State
	}
	if len(osr.Side) > 0 {
		params["side"] = osr.Side
	}
	if osr.Limit > 0 {
		params["limit"] = strconv.Itoa(osr.Limit)
	}
	if osr.RecvWindow > 0 {
		params["recv_window"] = strconv.Itoa(osr.RecvWindow)
	}
	res, err := ws.request("GET", "api/v1/orders", params, true, true)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from GetOrder:%s", err.Error()))
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ws.handleError(textRes)
	}

	var rawResult struct {
		Data []Order `json:"data"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.New(fmt.Sprintf("GetOrder Response unmarshal failed:%s", err.Error()))
	}
	var orders []*Order
	for _, v := range rawResult.Data {
		orders = append(orders, &v)
	}
	return orders, nil
}
func (ws *wonService) GetOrder(or OrderRequest) (*Order, error) {
	params := make(map[string]string)
	params["id"] = strconv.FormatInt(or.Id, 10)
	params["timestamp"] = strconv.FormatInt(or.Timestamp, 10)
	if or.RecvWindow > 0 {
		params["recv_window"] = strconv.Itoa(or.RecvWindow)
	}
	res, err := ws.request("GET", "api/v1/order", params, true, true)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from GetOrder:%s", err.Error()))
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ws.handleError(textRes)
	}

	var rawResult struct {
		Data Order `json:"data"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.New(fmt.Sprintf("GetOrder Response unmarshal failed:%s", err.Error()))
	}

	return &rawResult.Data, nil
}

func (ws *wonService) CancelOrder(cor CancelOrderRequest) error {
	params := make(map[string]string)
	params["id"] = strconv.FormatInt(cor.Id, 10)
	params["timestamp"] = strconv.FormatInt(cor.Timestamp, 10)
	if cor.RecvWindow > 0 {
		params["recv_window"] = strconv.Itoa(cor.RecvWindow)
	}
	res, err := ws.request("POST", "api/v1/order/cancel", params, true, true)
	if err != nil {
		return err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to read response from CancelOrder:%s", err.Error()))
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return ws.handleError(textRes)
	}

	var rawResult struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil || rawResult.Data != "success" {
		return errors.New(fmt.Sprintf("CancelOrder Response unmarshal failed:%s", err.Error()))
	}

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

func (ws *wonService) handleError(textRes []byte) error {
	err := &WonError{}
	level.Info(ws.Logger).Log("errorResponse", textRes)

	if err := json.Unmarshal(textRes, err); err != nil {
		return errors.New(fmt.Sprintf("error unmarshal failed:%s", err.Error()))
	}
	return err
}

func timeFromUnixMillTimestamp(ts int64) (time.Time, error) {
	return time.Unix(0, int64(ts)*int64(time.Millisecond)), nil
}
