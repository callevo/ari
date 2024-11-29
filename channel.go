package ari

import (
	"github.com/callevo/ari/channel"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/requests"
)

type ichannel struct {
	c *ARIClient
}

func (c *ichannel) Answer(key *key.Key) error {

	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelAnswer",
		Key:  key,
	})
}

func (c *ichannel) Get(key *key.Key) *channel.ChannelHandle {
	return nil
}

func (c *ichannel) Hangup(key *key.Key, reason string) error {
	return c.c.commandRequest((&requests.Request{
		Kind: "ChannelHangup",
		Key:  key,
		ChannelHangup: &requests.ChannelHangup{
			Reason: reason,
		},
	}))
}
