package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Service interface {
	Time() (time.Time, error)
	Depth() (DepthResult, error)
	Trades() (Trades, error)
	HistoricalTrades()(Trades, error)
	Account()(Account, error)
	TickerPrice()(TickerPrice, error)
	CreateOrder()(Order, error)
	GetOrders()([]*Order, error)
	GetOrder()(Order, error)
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
		return time.Time{}, errors.New(fmt.Sprintf("timeResponse unmarshal failed:%s", err.Error()))
	}
	t, err := timeFromUnixMillTimestamp(rawTime.Date.Time)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
	return time.Now(), nil
}

func(ws *wonService) Depth()(DepthResult, error){
	return DepthResult{}, nil
}
func(ws *wonService) Trades()(Trades, error){
	return  Trades{}, nil
}
func(ws *wonService) HistoricalTrades()(Trades, error){
	return  Trades{}, nil
}
func(ws *wonService) Account()(Account, error){
	return Account{}, nil
}
func(ws *wonService) TickerPrice()(TickerPrice, error){
	return TickerPrice{}, nil
}
func(ws *wonService) CreateOrder()(Order, error){
	return Order{}, nil
}
func(ws *wonService) GetOrders()([]*Order, error){
	return []*Order{}, nil
}
func(ws *wonService) GetOrder()(Order, error){
	return Order{}, nil
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

