package tests

import (
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