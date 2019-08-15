package pkg


import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Signer interface {
	Sign(payload []byte) string
}

type HmacSigner struct {
	Key []byte
}


func (h *HmacSigner) Sign(payload []byte) string {
	m := hmac.New(sha256.New, h.Key)
	m.Write(payload)
	return hex.EncodeToString(m.Sum(nil))
}
