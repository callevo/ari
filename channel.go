package ari

import (
	"time"

	"github.com/callevo/ari/arioptions"
	"github.com/callevo/ari/channel"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/requests"
)

type ichannel struct {
	c *ARIClient
}

func (c *ichannel) Create(ikey *key.Key, o requests.ChannelCreateRequest) (*channel.ChannelHandle, error) {
	k, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelCreate",
		Key:  ikey,
		ChannelCreate: &requests.ChannelCreate{
			ChannelCreateRequest: o,
		},
	})
	if err != nil {
		return nil, err
	}
	return channel.NewChannelHandle(k.New(key.ChannelKey, o.ChannelID), c, nil), nil
}

func (c *ichannel) Ring(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelRing",
		Key:  key,
	})
}

func (c *ichannel) StopRing(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelStopRing",
		Key:  key,
	})
}

func (c *ichannel) Answer(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelAnswer",
		Key:  key,
	})
}

func (c *ichannel) Busy(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelBusy",
		Key:  key,
	})
}

func (c *ichannel) Congestion(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelCongestion",
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

func (c *ichannel) Data(key *key.Key) (*channel.ChannelData, error) {
	data, err := c.c.dataRequest(&requests.Request{
		Kind: "ChannelData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Channel, nil
}

func (c *ichannel) Continue(key *key.Key, context string, extension string, priority int) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelContinue",
		Key:  key,
		ChannelContinue: &requests.ChannelContinue{
			Context:   context,
			Extension: extension,
			Priority:  priority,
		},
	})
}

func (c *ichannel) Dial(key *key.Key, caller string, timeout time.Duration) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelDial",
		Key:  key,
		ChannelDial: &requests.ChannelDial{
			Caller:  caller,
			Timeout: timeout,
		},
	})
}

func (c *ichannel) GetVariable(key *key.Key, name string) (string, error) {
	data, err := c.c.dataRequest(&requests.Request{
		Kind: "ChannelVariableGet",
		Key:  key,
		ChannelVariable: &requests.ChannelVariable{
			Name: name,
		},
	})
	if err != nil {
		return "", err
	}
	return data.Variable, nil
}

func (c *ichannel) SetVariable(key *key.Key, name, value string) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelVariableSet",
		Key:  key,
		ChannelVariable: &requests.ChannelVariable{
			Name:  name,
			Value: value,
		},
	})
}

func (c *ichannel) SendDTMF(key *key.Key, dtmf string, opts *arioptions.DTMFOptions) error {
	if opts == nil {
		opts = &arioptions.DTMFOptions{}
	}
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelSendDTMF",
		Key:  key,
		ChannelSendDTMF: &requests.ChannelSendDTMF{
			DTMF:    dtmf,
			Options: opts,
		},
	})
}

func (c *ichannel) Snoop(ikey *key.Key, snoopID string, opts *arioptions.SnoopOptions) (*channel.ChannelHandle, error) {
	k, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelSnoop",
		Key:  ikey,
		ChannelSnoop: &requests.ChannelSnoop{
			SnoopID: snoopID,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return channel.NewChannelHandle(k.New(key.ChannelKey, snoopID), c, nil), nil
}

func (c *ichannel) Hold(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelHold",
		Key:  key,
	})
}

func (c *ichannel) StopHold(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelStopHold",
		Key:  key,
	})
}

func (c *ichannel) Mute(key *key.Key, dir arioptions.Direction) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelMute",
		Key:  key,
		ChannelMute: &requests.ChannelMute{
			Direction: dir,
		},
	})
}

func (c *ichannel) Unmute(key *key.Key, dir arioptions.Direction) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelUnmute",
		Key:  key,
		ChannelMute: &requests.ChannelMute{
			Direction: dir,
		},
	})
}

func (c *ichannel) MOH(key *key.Key, moh string) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelMOH",
		Key:  key,
		ChannelMOH: &requests.ChannelMOH{
			Music: moh,
		},
	})
}

func (c *ichannel) StopMOH(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelStopMOH",
		Key:  key,
	})
}

func (c *ichannel) Silence(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelSilence",
		Key:  key,
	})
}

func (c *ichannel) StopSilence(key *key.Key) error {
	return c.c.commandRequest(&requests.Request{
		Kind: "ChannelStopSilence",
		Key:  key,
	})
}

func (c *ichannel) Originate(referenceKey *key.Key, o requests.OriginateRequest) (*channel.ChannelHandle, error) {
	k, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelOriginate",
		Key:  referenceKey,
		ChannelOriginate: &requests.ChannelOriginate{
			OriginateRequest: o,
		},
	})
	if err != nil {
		return nil, err
	}
	return channel.NewChannelHandle(k, c, nil), nil
}
