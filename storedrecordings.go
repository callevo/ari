package ari

import (
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/recordings"
	"github.com/callevo/ari/requests"
)

type iStoredRecording struct {
	c *ARIClient
}

func (s *iStoredRecording) List(filter *key.Key) ([]*key.Key, error) {
	return s.c.listRequest(&requests.Request{
		Kind: "RecordingStoredList",
		Key:  filter,
	})
}

func (s *iStoredRecording) Get(key *key.Key) *recordings.StoredRecordingHandle {
	k, err := s.c.getRequest(&requests.Request{
		Kind: "RecordingStoredGet",
		Key:  key,
	})
	if err != nil {
		logs.TLogger.Warn().Msgf("failed to get stored recording for handle %s", err)
		return recordings.NewStoredRecordingHandle(key, s, nil)
	}
	return recordings.NewStoredRecordingHandle(k, s, nil)
}

func (s *iStoredRecording) Data(key *key.Key) (*recordings.StoredRecordingData, error) {
	data, err := s.c.dataRequest(&requests.Request{
		Kind: "RecordingStoredData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.StoredRecording, nil
}

func (s *iStoredRecording) Copy(ikey *key.Key, dest string) (*recordings.StoredRecordingHandle, error) {
	h := recordings.NewStoredRecordingHandle(ikey.New(key.StoredRecordingKey, dest), s, nil)

	err := s.c.commandRequest(&requests.Request{
		Kind: "RecordingStoredCopy",
		Key:  ikey,
		RecordingStoredCopy: &requests.RecordingStoredCopy{
			Destination: dest,
		},
	})

	// NOTE: Always return the handle, even when we have an error
	return h, err
}

func (s *iStoredRecording) Delete(key *key.Key) error {
	return s.c.commandRequest(&requests.Request{
		Kind: "RecordingStoredDelete",
		Key:  key,
	})
}
