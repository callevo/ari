package recordings

import (
	"sync"
	"time"

	"github.com/callevo/ari/key"
)

// LiveRecording represents a communication path interacting with an Asterisk
// server for live recording resources
type LiveRecording interface {

	// Get gets the Recording by type
	Get(key *key.Key) *LiveRecordingHandle

	// Data gets the data for the live recording
	Data(key *key.Key) (*LiveRecordingData, error)

	// Stop stops the live recording
	Stop(key *key.Key) error

	// Pause pauses the live recording
	Pause(key *key.Key) error

	// Resume resumes the live recording
	Resume(key *key.Key) error

	// Mute mutes the live recording
	Mute(key *key.Key) error

	// Unmute unmutes the live recording
	Unmute(key *key.Key) error

	// Scrap Stops and deletes the current LiveRecording
	Scrap(key *key.Key) error

	// Stored returns the StoredRecording handle for this LiveRecording
	Stored(key *key.Key) *StoredRecordingHandle

	// Subscribe subscribes to events
	//Subscribe(key *Key, n ...string) Subscription
}

// LiveRecordingData is the data for a live recording
type LiveRecordingData struct {
	// Key is the cluster-unique identifier for this live recording
	Key *key.Key `json:"key"`

	Cause     string        `json:"cause,omitempty"`            // If failed, the cause of the failure
	Duration  time.Duration `json:"duration,omitempty"`         // Length of recording in seconds
	Format    string        `json:"format"`                     // Format of recording (wav, gsm, etc)
	Name      string        `json:"name"`                       // (base) name for the recording
	Silence   time.Duration `json:"silence_duration,omitempty"` // If silence was detected in the recording, the duration in seconds of that silence (requires that maxSilenceSeconds be non-zero)
	State     string        `json:"state"`                      // Current state of the recording
	Talking   time.Duration `json:"talking_duration,omitempty"` // Duration of talking, in seconds, that has been detected in the recording (requires that maxSilenceSeconds be non-zero)
	TargetURI string        `json:"target_uri"`                 // URI for the channel or bridge which is being recorded (TODO: figure out format for this)
}

// ID returns the identifier of the live recording
func (s *LiveRecordingData) ID() string {
	return s.Name
}

// NewLiveRecordingHandle creates a new live recording handle
func NewLiveRecordingHandle(ikey *key.Key, r LiveRecording, exec func(*LiveRecordingHandle) (err error)) *LiveRecordingHandle {
	return &LiveRecordingHandle{
		key:  ikey,
		r:    r,
		exec: exec,
	}
}

// A LiveRecordingHandle is a reference to a live recording that can be operated on
type LiveRecordingHandle struct {
	key      *key.Key
	r        LiveRecording
	exec     func(*LiveRecordingHandle) (err error)
	executed bool

	mu sync.Mutex
}

// ID returns the identifier of the live recording
func (h *LiveRecordingHandle) ID() string {
	return h.key.ID
}

// Key returns the key of the live recording
func (h *LiveRecordingHandle) Key() *key.Key {
	return h.key
}

// Data gets the data for the live recording
func (h *LiveRecordingHandle) Data() (*LiveRecordingData, error) {
	return h.r.Data(h.key)
}

// Stop stops and saves the recording
func (h *LiveRecordingHandle) Stop() error {
	return h.r.Stop(h.key)
}

// Scrap stops and deletes the recording
func (h *LiveRecordingHandle) Scrap() error {
	return h.r.Scrap(h.key)
}

// Resume resumes the recording
func (h *LiveRecordingHandle) Resume() error {
	return h.r.Resume(h.key)
}

// Pause pauses the recording
func (h *LiveRecordingHandle) Pause() error {
	return h.r.Pause(h.key)
}

// Mute mutes the recording
func (h *LiveRecordingHandle) Mute() error {
	return h.r.Mute(h.key)
}

// Unmute mutes the recording
func (h *LiveRecordingHandle) Unmute() error {
	return h.r.Unmute(h.key)
}

// Stored returns the StoredRecordingHandle for this LiveRecordingHandle
func (h *LiveRecordingHandle) Stored() *StoredRecordingHandle {
	return h.r.Stored(h.key)
}

// Exec executes any staged operations attached to the `LiveRecordingHandle`
func (h *LiveRecordingHandle) Exec() (err error) {
	h.mu.Lock()

	if !h.executed {
		h.executed = true
		if h.exec != nil {
			err = h.exec(h)
			h.exec = nil
		}
	}

	h.mu.Unlock()

	return
}
