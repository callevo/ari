package ari

import (
	"github.com/callevo/ari/bridge"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/requests"
)

type ibridge struct {
	c *ARIClient
}

func (b *ibridge) AddChannel(key *key.Key, channelID string) error {
	return b.AddChannelWithOptions(key, channelID, nil)
}

func (b *ibridge) Create(key *key.Key, btype, name string) (*bridge.BridgeHandle, error) {
	k, err := b.c.createRequest(&requests.Request{
		Kind: "BridgeCreate",
		Key:  key,
		BridgeCreate: &requests.BridgeCreate{
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
		BridgeCreate: &requests.BridgeCreate{
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
		BridgeRemoveChannel: &requests.BridgeRemoveChannel{
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

func (b *ibridge) StopMOH(key *key.Key) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeStopMOH",
		Key:  key,
	})
}

func (b *ibridge) Data(key *key.Key) (*bridge.BridgeData, error) {
	resp, err := b.c.dataRequest(&requests.Request{
		Kind: "BridgeData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Bridge, nil
}

func (b *ibridge) Get(key *key.Key) *bridge.BridgeHandle {
	k, err := b.c.getRequest(&requests.Request{
		Kind: "BridgeGet",
		Key:  key,
	})
	if err != nil {
		logs.TLogger.Error().Msgf("failed to get bridge for handle %s", err.Error())

		return bridge.NewBridgeHandle(key, b, nil)
	}
	return bridge.NewBridgeHandle(k, b, nil)
}

func (b *ibridge) VideoSource(key *key.Key, channelID string) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeVideoSource",
		Key:  key,
		BridgeVideoSource: &requests.BridgeVideoSource{
			Channel: channelID,
		},
	})
}

func (b *ibridge) VideoSourceDelete(key *key.Key) error {
	return b.c.commandRequest(&requests.Request{
		Kind: "BridgeVideoSourceDelete",
		Key:  key,
	})
}
