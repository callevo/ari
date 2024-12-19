package channel

import (
	"strings"
	"time"

	"github.com/callevo/ari/arioptions"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/play"
	"github.com/callevo/ari/requests"
)

type Channel interface {
	// Get returns a handle to a channel for further interaction
	Get(key *key.Key) *ChannelHandle

	// GetVariable retrieves the value of a channel variable
	GetVariable(*key.Key, string) (string, error)

	// List lists the channels in asterisk, optionally using the key for filtering
	//List(*key.Key) ([]*key.Key, error)

	// Originate creates a new channel, returning a handle to it or an error, if
	// the creation failed.
	// The Key should be that of the linked channel, if one exists, so that the
	// Node can be matches to it.
	Originate(*key.Key, requests.OriginateRequest) (*ChannelHandle, error)

	// StageOriginate creates a new Originate, created when the `Exec` method
	// on `ChannelHandle` is invoked.
	// The Key should be that of the linked channel, if one exists, so that the
	// Node can be matches to it.
	//StageOriginate(*key.Key, requests.OriginateRequest) (*ChannelHandle, error)

	// Create creates a new channel, returning a handle to it or an
	// error, if the creation failed. Create is already Staged via `Dial`.
	// The Key should be that of the linked channel, if one exists, so that the
	// Node can be matches to it.
	Create(*key.Key, requests.ChannelCreateRequest) (*ChannelHandle, error)

	// Data returns the channel data for a given channel
	Data(key *key.Key) (*ChannelData, error)

	// Continue tells Asterisk to return a channel to the dialplan
	Continue(key *key.Key, context, extension string, priority int) error

	// Busy hangs up the channel with the "busy" cause code
	Busy(key *key.Key) error

	// Congestion hangs up the channel with the "congestion" cause code
	Congestion(key *key.Key) error

	// Answer answers the channel
	Answer(key *key.Key) error

	// Hangup hangs up the given channel
	Hangup(key *key.Key, reason string) error

	// Ring indicates ringing to the channel
	Ring(key *key.Key) error

	// StopRing stops ringing on the channel
	StopRing(key *key.Key) error

	// SendDTMF sends DTMF to the channel
	SendDTMF(key *key.Key, dtmf string, opts *arioptions.DTMFOptions) error

	// Hold puts the channel on hold
	Hold(key *key.Key) error

	// StopHold retrieves the channel from hold
	StopHold(key *key.Key) error

	// Mute mutes a channel in the given direction (in,out,both)
	Mute(key *key.Key, dir arioptions.Direction) error

	// Unmute unmutes a channel in the given direction (in,out,both)
	Unmute(key *key.Key, dir arioptions.Direction) error

	// MOH plays music on hold
	MOH(key *key.Key, moh string) error

	// SetVariable sets a channel variable
	SetVariable(key *key.Key, name, value string) error

	// StopMOH stops music on hold
	StopMOH(key *key.Key) error

	// Silence plays silence to the channel
	Silence(key *key.Key) error

	// StopSilence stops the silence on the channel
	StopSilence(key *key.Key) error

	// Play plays the media URI to the channel
	Play(key *key.Key, playbackID string, mediaURI string) (*play.PlaybackHandle, error)

	// StagePlay stages a `Play` operation and returns the `PlaybackHandle`
	// for invoking it.
	//StagePlay(key *key.Key, playbackID string, mediaURI string) (*PlaybackHandle, error)

	// Record records the channel
	//Record(key *Key, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// StageRecord stages a `Record` operation and returns the `PlaybackHandle`
	// for invoking it.
	//StageRecord(key *Key, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// Dial dials a created channel
	Dial(key *key.Key, caller string, timeout time.Duration) error

	// Snoop spies on a specific channel, creating a new snooping channel
	Snoop(key *key.Key, snoopID string, opts *arioptions.SnoopOptions) (*ChannelHandle, error)

	// StageSnoop creates a new `ChannelHandle`, when `Exec`ed, snoops on the given channel ID and
	// creates a new snooping channel.
	//StageSnoop(key *key.Key, snoopID string, opts *SnoopOptions) (*ChannelHandle, error)

	// StageExternalMedia creates a new non-telephony external media channel,
	// when `Exec`ed, by which audio may be sent or received.  The stage version
	// of this command will not actually communicate with Asterisk until Exec is
	// called on the returned ExternalMedia channel.
	StageExternalMedia(key *key.Key, opts arioptions.ExternalMediaOptions) (*ChannelHandle, error)

	// ExternalMedia creates a new non-telephony external media channel by which audio may be sent or received
	ExternalMedia(key *key.Key, opts arioptions.ExternalMediaOptions) (*ChannelHandle, error)

	// UserEvent Sends user-event to AMI channel subscribers
	//UserEvent(key *Key, ue *ChannelUserevent) error
}

type CallerInfo struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type ChannelData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	Protocol     string `json:"protocol"`
	Caller       CallerInfo
	Connected    CallerInfo
	Accountcode  string `json:"accountcode"`
	Dialplan     DialplanInfo
	Creationtime string            `json:"creationtime"`
	Language     string            `json:"language"`
	ChannelVars  map[string]string `json:"channelvars"`
}

func (c *ChannelData) GetID() string {
	return c.ID
}

type DialplanInfo struct {
	Context  string `json:"context"`
	Exten    string `json:"exten"`
	Priority int    `json:"priority"`
	AppName  string `json:"app_name"`
	AppData  string `json:"app_data"`
}

// ChannelHandle provides a wrapper on the Channel interface for operations on a particular channel ID.
type ChannelHandle struct {
	key *key.Key

	c Channel

	callback func(ch *ChannelHandle) error

	executed bool
}

// Exec executes any staged channel operations attached to this handle.
func (ch *ChannelHandle) Exec() (err error) {
	if !ch.executed {
		ch.executed = true
		if ch.callback != nil {
			err = ch.callback(ch)
			ch.callback = nil
		}
	}

	return err
}

// NewChannelHandle returns a handle to the given ARI channel
func NewChannelHandle(key *key.Key, c Channel, exec func(ch *ChannelHandle) error) *ChannelHandle {
	return &ChannelHandle{
		key:      key,
		c:        c,
		callback: exec,
	}
}

func (ch *ChannelHandle) ID() string {
	return ch.key.ID
}

// Key returns the key for the channel handle
func (ch *ChannelHandle) Key() *key.Key {
	return ch.key
}

// Data returns the channel's data
func (ch *ChannelHandle) Data(key *key.Key) (*ChannelData, error) {
	return ch.c.Data(key)
}

// Continue tells Asterisk to return the channel to the dialplan
func (ch *ChannelHandle) Continue(context, extension string, priority int) error {
	return ch.c.Continue(ch.key, context, extension, priority)
}

// Play initiates playback of the specified media uri
// to the channel, returning the Playback handle
func (ch *ChannelHandle) Play(id string, mediaURI string) (ph *play.PlaybackHandle, err error) {
	return ch.c.Play(ch.key, id, mediaURI)
}

// Record records the channel to the given filename
//func (ch *ChannelHandle) Record(name string, opts *RecordingOptions) (*LiveRecordingHandle, error) {
//	return nil, nil//
//}

//---
// Hangup Operations
//---

// Busy hangs up the channel with the "busy" cause code
func (ch *ChannelHandle) Busy() error {
	return ch.c.Busy(ch.key)
}

// Congestion hangs up the channel with the congestion cause code
func (ch *ChannelHandle) Congestion() error {
	return ch.c.Congestion(ch.key)
}

// Hangup hangs up the channel with the normal cause code
func (ch *ChannelHandle) Hangup() error {
	return ch.c.Hangup(ch.key, "normal")
}

// Answer operations
// --

// Answer answers the channel
func (ch *ChannelHandle) Answer() error {
	return ch.c.Answer(ch.key)
}

// IsAnswered checks the current state of the channel to see if it is "Up"
func (ch *ChannelHandle) IsAnswered() (bool, error) {
	updated, err := ch.Data(ch.key)
	if err != nil {
		return false, err
	}

	return strings.ToLower(updated.State) == "up", nil
}

// Ring Operations
// --

// Ring indicates ringing to the channel
func (ch *ChannelHandle) StopRing() error {
	return ch.c.StopRing(ch.key)
}

// StopRing stops ringing on the channel
func (ch *ChannelHandle) Ring() error {
	return ch.c.Ring(ch.key)
}

// Mute operations
// --

// Mute mutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Mute(dir arioptions.Direction) (err error) {
	if dir == "" {
		dir = arioptions.DirectionIn
	}

	return ch.c.Mute(ch.key, dir)
}

// Unmute unmutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Unmute(dir arioptions.Direction) (err error) {
	if dir == "" {
		dir = arioptions.DirectionIn
	}

	return ch.c.Unmute(ch.key, dir)
}

// Hold operations
// --

// Hold puts the channel on hold
func (ch *ChannelHandle) Hold() error {
	return ch.c.Hold(ch.key)
}

// StopHold retrieves the channel from hold
func (ch *ChannelHandle) StopHold() error {
	return ch.c.StopHold(ch.key)
}

// Music on hold operations
// --

// MOH plays music on hold of the given class
// to the channel
func (ch *ChannelHandle) MOH(mohClass string) error {
	return ch.c.MOH(ch.key, mohClass)
}

// StopMOH stops playing of music on hold to the channel
func (ch *ChannelHandle) StopMOH() error {
	return ch.c.StopMOH(ch.key)
}

// ----

// GetVariable returns the value of a channel variable
func (ch *ChannelHandle) GetVariable(name string) (string, error) {
	return ch.c.GetVariable(ch.key, name)
}

// SetVariable sets the value of a channel variable
func (ch *ChannelHandle) SetVariable(name, value string) error {
	return ch.c.SetVariable(ch.key, name, value)
}

// --
// Misc
// --

// Originate creates (and dials) a new channel using the present channel as its Originator.
func (ch *ChannelHandle) Originate(req requests.OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return ch.c.Originate(ch.key, req)
}

// StageOriginate stages an originate (channel creation and dial) to be Executed later.
func (ch *ChannelHandle) StageOriginate(req arioptions.OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return nil, nil
}

// Create creates (but does not dial) a new channel, using the present channel as its Originator.
func (ch *ChannelHandle) Create(req requests.ChannelCreateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return ch.c.Create(ch.key, req)
}

// Dial dials a created channel.  `caller` is the optional
// channel ID of the calling party (if there is one).  Timeout
// is the length of time to wait before the dial is answered
// before aborting.
func (ch *ChannelHandle) Dial(caller string, timeout time.Duration) error {
	return ch.c.Dial(ch.key, caller, timeout)
}

// Snoop spies on a specific channel, creating a new snooping channel placed into the given app
func (ch *ChannelHandle) Snoop(snoopID string, opts *arioptions.SnoopOptions) (*ChannelHandle, error) {
	return ch.c.Snoop(ch.key, snoopID, opts)
}

// StageSnoop stages a `Snoop` operation
func (ch *ChannelHandle) StageSnoop(snoopID string, opts *arioptions.SnoopOptions) (*ChannelHandle, error) {
	return nil, nil
}

// StageExternalMedia creates a new non-telephony external media channel,
// when `Exec`ed, by which audio may be sent or received.  The stage version
// of this command will not actually communicate with Asterisk until Exec is
// called on the returned ExternalMedia channel.
func (ch *ChannelHandle) StageExternalMedia(opts arioptions.ExternalMediaOptions) (*ChannelHandle, error) {
	return ch.c.StageExternalMedia(ch.key, opts)
}

// ExternalMedia creates a new non-telephony external media channel by which audio may be sent or received
func (ch *ChannelHandle) ExternalMedia(opts arioptions.ExternalMediaOptions) (*ChannelHandle, error) {
	return ch.c.ExternalMedia(ch.key, opts)
}

// Silence operations
// --

// Silence plays silence to the channel
func (ch *ChannelHandle) Silence() error {
	return ch.c.Silence(ch.key)
}

// StopSilence stops silence to the channel
func (ch *ChannelHandle) StopSilence() error {
	return ch.c.StopSilence(ch.key)
}

// DTMF
// --

// SendDTMF sends the DTMF information to the server
func (ch *ChannelHandle) SendDTMF(dtmf string, opts *arioptions.DTMFOptions) error {
	return ch.c.SendDTMF(ch.key, dtmf, opts)
}

// UserEvent sends user-event to AMI channel subscribers
//func (ch *ChannelHandle) UserEvent(key *Key, ue *ChannelUserevent) error {
//	return nil
//}

//Code taken from ari-proxy. NOt verbatim, but with enough similarities to make this work .
//Key was a copy from ari-proxy
