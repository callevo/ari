package ari

import (
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/play"
	"github.com/callevo/ari/requests"
)

type playback struct {
	c *ARIClient
}

func (p *playback) Get(key *key.Key) *play.PlaybackHandle {
	k, err := p.c.getRequest(&requests.Request{
		Kind: "PlaybackGet",
		Key:  key,
	})
	if err != nil {
		logs.TLogger.Warn().Msgf("failed to get playback for handle %s", err)
		return play.NewPlaybackHandle(key, p, nil)
	}

	return play.NewPlaybackHandle(k, p, nil)
}

func (p *playback) Data(key *key.Key) (*play.PlaybackData, error) {
	data, err := p.c.dataRequest(&requests.Request{
		Kind: "PlaybackData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Playback, nil
}

func (p *playback) Control(key *key.Key, op string) error {
	return p.c.commandRequest(&requests.Request{
		Kind: "PlaybackControl",
		Key:  key,
		PlaybackControl: &requests.PlaybackControl{
			Command: op,
		},
	})
}

func (p *playback) Stop(key *key.Key) error {
	return p.c.commandRequest(&requests.Request{
		Kind: "PlaybackStop",
		Key:  key,
	})
}
