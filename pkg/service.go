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
	Trades(TradeRequest) ([]*Trade, error)
	HistoricalTrades(TradeRequest) ([]*HistoryTrade, error)
	Account(AccountRequest) (*Account, error)
	TickerPrice(TickerPriceRequest) (*TickerPrice, error)
	CreateOrder() (Order, error)
	GetOrders() ([]*Order, error)
	GetOrder() (Order, error)
	CancelOrder() error
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
func (ws *wonService) Trades(tr TradeRequest) ([]*Trade, error) {
	params := make(map[string]string)
	params["market"] = tr.Market
	if tr.Limit > 0 {
		params["limit"] = strconv.Itoa(tr.Limit)
	}

	if tr.FromId > 0 {
		params["from_id"] = strconv.Itoa(tr.FromId)
	}

	res, err := ws.request("GET", "api/v1/trades", params, false, false)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from Depth:%s", err.Error()))
	}
	defer res.Body.Close()

	type result struct {
		Id       int    `json:"id"`
		Price    string `json:"price"`
		Quantity string `json:"qty"`
		CreateAt int    `json:"time"`
	}

	var rawDepth struct {
		Data []result `json:"data"`
	}

	if err := json.Unmarshal(textRes, &rawDepth); err != nil {
		return nil, errors.New(fmt.Sprintf("Trades Response unmarshal Trades:%s", err.Error()))
	}
	var trades []*Trade
	for _, v := range rawDepth.Data {
		trades = append(trades, &Trade{Id: v.Id, Price: v.Price, Quantity: v.Quantity, CreateAt: v.CreateAt})
	}

	return trades, nil
}
func (ws *wonService) HistoricalTrades(tr TradeRequest) ([]*HistoryTrade, error) {
	params := make(map[string]string)
	params["market"] = tr.Market
	if tr.Limit > 0 {
		params["limit"] = strconv.Itoa(tr.Limit)
	}

	if tr.FromId > 0 {
		params["from_id"] = strconv.Itoa(tr.FromId)
	}

	res, err := ws.request("GET", "api/v1/history/trades", params, true, true)
	if err != nil {
		return nil, err
	}

	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response from HistoricalTrades:%s", err.Error()))
	}
	defer res.Body.Close()

	type result struct {
		Id       int    `json:"id"`
		Price    string `json:"price"`
		Side     string `json:"side"`
		Quantity string `json:"qty"`
		CreateAt int    `json:"time"`
	}

	var rawDepth struct {
		Data []result `json:"data"`
	}

	if err := json.Unmarshal(textRes, &rawDepth); err != nil {
		return nil, errors.New(fmt.Sprintf("HistoricalTrades Response unmarshal HistoricalTrades:%s", err.Error()))
	}
	var trades []*HistoryTrade
	for _, v := range rawDepth.Data {
		trades = append(trades, &HistoryTrade{Id: v.Id, Price: v.Price, Side: v.Side, Quantity: v.Quantity, CreateAt: v.CreateAt})
	}

	return trades, nil
}
func (ws *wonService) Account(ar AccountRequest) (*Account, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.Itoa(ar.Timestamp)
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
		Precision    int `json:"precision"`
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

func (ws *wonService) CreateOrder() (Order, error) {
	return Order{}, nil
}
func (ws *wonService) GetOrders() ([]*Order, error) {
	return []*Order{}, nil
}
func (ws *wonService) GetOrder() (Order, error) {
	return Order{}, nil
}
func (ws *wonService) CancelOrder() error {
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
