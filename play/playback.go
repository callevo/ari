package play

import "github.com/callevo/ari/key"

// Playback represents a communication path for interacting
// with an Asterisk server for playback resources
type Playback interface {

	// Get gets the handle to the given playback ID
	Get(key *key.Key) *PlaybackHandle

	// Data gets the playback data
	Data(key *key.Key) (*PlaybackData, error)

	// Control performs the given operation on the current playback.  Available operations are:
	//   - restart
	//   - pause
	//   - unpause
	//   - reverse
	//   - forward
	Control(key *key.Key, op string) error

	// Stop stops the playback
	Stop(key *key.Key) error
}

// A Player is an entity which can play an audio URI
type Player interface {
	// Play plays the audio using the given playback ID and media URI
	Play(string, string) (*PlaybackHandle, error)

	// StagePlay stages a `Play` operation
	StagePlay(string, string) (*PlaybackHandle, error)
}

// PlaybackData represents the state of a playback
type PlaybackData struct {
	// Key is the cluster-unique identifier for this playback
	Key *key.Key `json:"key"`

	ID        string `json:"id"` // Unique ID for this playback session
	Language  string `json:"language,omitempty"`
	MediaURI  string `json:"media_uri"`  // URI for the media which is to be played
	State     string `json:"state"`      // State of the playback operation
	TargetURI string `json:"target_uri"` // URI of the channel or bridge on which the media should be played (follows format of 'type':'name')
}

// PlaybackHandle is the handle for performing playback operations
type PlaybackHandle struct {
	key      *key.Key
	p        Playback
	exec     func(pb *PlaybackHandle) error
	executed bool
}

// NewPlaybackHandle builds a handle to the playback id
func NewPlaybackHandle(key *key.Key, pb Playback, exec func(pb *PlaybackHandle) error) *PlaybackHandle {
	return &PlaybackHandle{
		key:  key,
		p:    pb,
		exec: exec,
	}
}

// ID returns the identifier for the playback
func (ph *PlaybackHandle) ID() string {
	return ph.key.ID
}

// Key returns the Key for the playback
func (ph *PlaybackHandle) Key() *key.Key {
	return ph.key
}

// Data gets the playback data
func (ph *PlaybackHandle) Data() (*PlaybackData, error) {
	return ph.p.Data(ph.key)
}

// Control performs the given operation
func (ph *PlaybackHandle) Control(op string) error {
	return ph.p.Control(ph.key, op)
}

// Stop stops the playback
func (ph *PlaybackHandle) Stop() error {
	return ph.p.Stop(ph.key)
}

// Exec executes any staged operations
func (ph *PlaybackHandle) Exec() (err error) {
	if !ph.executed {
		ph.executed = true
		if ph.exec != nil {
			err = ph.exec(ph)
			ph.exec = nil
		}
	}

	return
}
