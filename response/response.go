package response

import (
	"errors"

	"github.com/callevo/ari/asterisk"
	"github.com/callevo/ari/bridge"
	"github.com/callevo/ari/channel"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/play"
	"github.com/callevo/ari/recordings"
)

// ErrNotFound indicates that the operation did not return a result
var ErrNotFound = errors.New("Not found")

type Response struct {
	// Error is the error encountered
	Error string `json:"error"`

	Code int `json:"code,omitempty"`

	// Data is the returned entity data, if applicable
	Data *EntityData `json:"data,omitempty"`

	// Key is the key of the returned entity, if applicable
	Key *key.Key `json:"key,omitempty"`

	// Keys is the list of keys of any matching entities, if applicable
	Keys []*key.Key `json:"keys,omitempty"`
}

type EntityData struct {
	Channel         *channel.ChannelData            `json:"channel,omitempty"`
	Asterisk        *asterisk.AsteriskInfo          `json:"asterisk,omitempty"`
	Bridge          *bridge.BridgeData              `json:"bridge,omitempty"`
	Playback        *play.PlaybackData              `json:"playback,omitempty"`
	StoredRecording *recordings.StoredRecordingData `json:"stored_recording,omitempty"`
	LiveRecording   *recordings.LiveRecordingData   `json:"live_recording,omitempty"`
	ChannelList     []*channel.ChannelData          `json:"channellist,omitempty"`

	Variable string `json:"variable,omitempty"`
}

// Err returns an error from the Response.  If the response's Error is empty, a nil error is returned.  Otherwise, the error will be filled with the value of response.Error.
func (e *Response) Err() error {
	if e == nil {
		return nil
	}
	if e.Error != "" {
		return errors.New(e.Error)
	}
	return nil
}

// IsNotFound indicates that the retuned error response was a Not Found error response
func (e *Response) IsNotFound() bool {
	return e.Error == "Not found"
}

// NewErrorResponse wraps an error as an ErrorResponse
func NewErrorResponse(err error) *Response {
	if err == nil {
		return &Response{}
	}
	return &Response{Error: err.Error()}
}
