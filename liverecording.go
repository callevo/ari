package ari

import (
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/recordings"
	"github.com/callevo/ari/requests"
)

type iLifeRecording struct {
	c *ARIClient
}

func (l *iLifeRecording) Get(key *key.Key) *recordings.LiveRecordingHandle {
	k, err := l.c.getRequest(&requests.Request{
		Kind: "RecordingLiveGet",
		Key:  key,
	})
	if err != nil {
		logs.TLogger.Warn().Msgf("failed to get liveRecording for handle %s", err)
		return recordings.NewLiveRecordingHandle(key, l, nil)
	}
	return recordings.NewLiveRecordingHandle(k, l, nil)
}

func (l *iLifeRecording) Data(key *key.Key) (*recordings.LiveRecordingData, error) {
	data, err := l.c.dataRequest(&requests.Request{
		Kind: "RecordingLiveData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.LiveRecording, nil
}

func (l *iLifeRecording) Stop(key *key.Key) error {
	return l.c.commandRequest(&requests.Request{
		Kind: "RecordingLiveStop",
		Key:  key,
	})
}

func (l *iLifeRecording) Pause(key *key.Key) error {
	return l.c.commandRequest(&requests.Request{
		Kind: "RecordingLivePause",
		Key:  key,
	})
}

func (l *iLifeRecording) Resume(key *key.Key) error {
	return l.c.commandRequest(&requests.Request{
		Kind: "RecordingLiveResume",
		Key:  key,
	})
}

func (l *iLifeRecording) Mute(key *key.Key) error {
	return l.c.commandRequest(&requests.Request{
		Kind: "RecordingLiveMute",
		Key:  key,
	})
}

func (l *iLifeRecording) Unmute(key *key.Key) error {
	return l.c.commandRequest(&requests.Request{
		Kind: "RecordingLiveUnmute",
		Key:  key,
	})
}

func (l *iLifeRecording) Scrap(ikey *key.Key) error {
	return l.c.commandRequest(&requests.Request{
		Kind: "RecordingLiveScrap",
		Key:  ikey,
	})
}

func (l *iLifeRecording) Stored(ikey *key.Key) *recordings.StoredRecordingHandle {
	return recordings.NewStoredRecordingHandle(ikey.New(key.StoredRecordingKey, ikey.ID), l.c.StoredRecording(), nil)
}
