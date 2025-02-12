package ari

import (
	"time"

	"github.com/callevo/ari/arioptions"
	"github.com/callevo/ari/channel"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/play"
	"github.com/callevo/ari/recordings"
	"github.com/callevo/ari/requests"
	"github.com/callevo/ari/rid"
)

type ichannel struct {
	c *ARIClient
}

func (c *ichannel) List(filter *key.Key) ([]*key.Key, error) {
	return c.c.listRequest(&requests.Request{
		Kind: "ChannelList",
		Key:  filter,
	})
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
	k, err := c.c.getRequest(&requests.Request{
		Kind: "ChannelGet",
		Key:  key,
	})
	if err != nil {
		logs.TLogger.Warn().Msgf("failed to make data request for channel %s", err)
		return channel.NewChannelHandle(key, c, nil)
	}
	return channel.NewChannelHandle(k, c, nil)
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

	if opts.App == "" {
		opts.App = c.c.ApplicationName()
	}
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

func (c *ichannel) Play(ikey *key.Key, playbackID string, mediaURI string) (*play.PlaybackHandle, error) {
	if playbackID == "" {
		playbackID = rid.New(rid.Playback)
	}

	k, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelPlay",
		Key:  ikey,
		ChannelPlay: &requests.ChannelPlay{
			PlaybackID: playbackID,
			MediaURI:   mediaURI,
		},
	})
	if err != nil {
		return nil, err
	}
	return play.NewPlaybackHandle(k.New(key.PlaybackKey, playbackID), c.c.Playback(), nil), nil
}

func (c *ichannel) Record(ikey *key.Key, name string, opts *arioptions.RecordingOptions) (*recordings.LiveRecordingHandle, error) {
	rb, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelRecord",
		Key:  ikey,
		ChannelRecord: &requests.ChannelRecord{
			Name:    name,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return recordings.NewLiveRecordingHandle(rb.New(key.LiveRecordingKey, name), c.c.LiveRecording(), nil), nil
}

func (c *ichannel) ExternalMedia(referenceKey *key.Key, opts arioptions.ExternalMediaOptions) (*channel.ChannelHandle, error) {
	if opts.ChannelID == "" {
		opts.ChannelID = rid.New(rid.Channel)
	}
	k, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelExternalMedia",
		Key:  referenceKey,
		ChannelExternalMedia: &requests.ChannelExternalMedia{
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return channel.NewChannelHandle(k, c, nil), nil
}

func (c *ichannel) StageExternalMedia(referenceKey *key.Key, opts arioptions.ExternalMediaOptions) (*channel.ChannelHandle, error) {
	if opts.ChannelID == "" {
		opts.ChannelID = rid.New(rid.Channel)
	}

	// We go ahead an call the createRequest on the server so that we lock in an
	// Asterisk box at the time of staging even though this staging call will
	// never actually be used.
	k, err := c.c.createRequest(&requests.Request{
		Kind: "ChannelStageOriginate",
		Key:  referenceKey,
		ChannelExternalMedia: &requests.ChannelExternalMedia{
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return channel.NewChannelHandle(k, c, func(h *channel.ChannelHandle) error {
		_, err := c.ExternalMedia(k, opts)
		return err
	}), nil
}
