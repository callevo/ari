package ari

import (
	"github.com/callevo/ari/asterisk"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/requests"
)

type iasterisk struct {
	c *ARIClient
}

func (a *iasterisk) Info(key *key.Key) (*asterisk.AsteriskInfo, error) {
	resp, err := a.c.dataRequest(&requests.Request{
		Kind: "AsteriskInfo",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Asterisk, nil
}
