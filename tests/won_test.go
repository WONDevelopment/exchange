package tests

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/xiangxian/exchange"
	"github.com/xiangxian/exchange/pkg"
	"testing"
	"time"
)

func TestTime(t *testing.T){
	signer:= &pkg.HmacSigner{}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	won := exchange.NewWon(server)
	m, err:=won.Time()
	t.Logf(m.String())
	assert.Equal(t, nil ,err)
	assert.Equal(t, true, m.UnixNano() < time.Now().UnixNano())
}


func TestDepth(t *testing.T){
	signer:= &pkg.HmacSigner{}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	won := exchange.NewWon(server)
	r, err:=won.Depth(pkg.DepthRequest{Market:"wonbtc",Limit:10})
	t.Logf(fmt.Sprintf("DepthResult:%v", r))
	assert.Equal(t, nil ,err)
}

func TestTrades(t *testing.T){
	signer:= &pkg.HmacSigner{}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	won := exchange.NewWon(server)
	r, err:=won.Trades(pkg.TradeRequest{Market:"wonbtc",Limit:10})
	t.Logf(fmt.Sprintf("TradeResult:%v", r))
	assert.Equal(t, nil ,err)
}


func TestHistoryTrades(t *testing.T){
	signer:= &pkg.HmacSigner{}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	won := exchange.NewWon(server)
	r, err:=won.HistoricalTrades(pkg.TradeRequest{Market:"wonbtc",Limit:10})
	t.Logf(fmt.Sprintf("HistoricalTrades:%v", r))
	assert.Equal(t, nil ,err)
}

func TestTickerPrice(t *testing.T){
	signer:= &pkg.HmacSigner{}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	won := exchange.NewWon(server)
	r, err:=won.TickerPrice(pkg.TickerPriceRequest{Market:"wonbtc"})
	t.Logf(fmt.Sprintf("TradeResult:%v", r))
	assert.Equal(t, nil ,err)
}

func TestAccount(t *testing.T){
	signer:= &pkg.HmacSigner{}
	server := pkg.NewWonService(
		wonHost,
		"",
		signer,
		nil,
		nil)
	won := exchange.NewWon(server)
	r, err:=won.Account(pkg.AccountRequest{Timestamp:int(time.Now().Unix() * 100), RecvWindow:5000})
	t.Logf(fmt.Sprintf("AccountResult:%v", r))
	assert.Equal(t, nil ,err)
}

