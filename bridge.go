package ari

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/callevo/ari/bridge"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/requests"
)

type ibridge struct {
	c *ARIClient
}

func (c *ibridge) AddChannel(key *key.Key, channelID string) error {
	return b.AddChannelWithOptions(key, channelID, nil)
}

func (b *ibridge) Create(key *key.Key, btype, name string) (*bridge.BridgeHandle, error) {
	k, err := b.c.createRequest(&request.Request{
		Kind: "BridgeCreate",
		Key:  key,
		BridgeCreate: &request.BridgeCreate{
			Type: btype,
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}
	return bridge.NewBridgeHandle(k, b, nil), nil
}

func (b *ibridge) StageCreate(key *key.Key, btype, name string) (*bridge.BridgeHandle, error) {
	k, err := b.c.createRequest(&requests.Request{
		Kind: "BridgeStageCreate",
		Key:  key,
		BridgeCreate: &request.BridgeCreate{
			Type: btype,
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}
	return bridge.NewBridgeHandle(k, b, func(h *bridge.BridgeHandle) error {
		_, err := b.Create(k, btype, name)
		return err
	}), nil
}

func (b *ibridge) AddChannelWithOptions(key *key.Key, channelID string, options *bridge.BridgeAddChannelOptions) error {
	if options == nil {
		options = new(bridge.BridgeAddChannelOptions)
	}

	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeAddChannel",
		Key:  key,
		BridgeAddChannel: &requests.BridgeAddChannel{
			Channel:    channelID,
			AbsorbDTMF: options.AbsorbDTMF,
			Mute:       options.Mute,
			Role:       options.Role,
		},
	})
}

func (b *ibridge) RemoveChannel(key *key.Key, channelID string) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeRemoveChannel",
		Key:  key,
		BridgeRemoveChannel: &request.BridgeRemoveChannel{
			Channel: channelID,
		},
	})
}

func (b *ibridge) Delete(key *key.Key) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeDelete",
		Key:  key,
	})
}

func (b *ibridge) MOH(key *key.Key, class string) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeMOH",
		Key:  key,
		BridgeMOH: &requests.BridgeMOH{
			Class: class,
		},
	})
}

func (b *bridge) StopMOH(key *key.Key) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeStopMOH",
		Key:  key,
	})
}
