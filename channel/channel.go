package channel

import (
	"strings"
	"time"

	"github.com/callevo/ari/arioptions"
	"github.com/callevo/ari/key"
)

type Channel interface {
	Get(key *key.Key) *ChannelHandle
	Answer(key *key.Key) error
	Hangup(key *key.Key, reason string) error
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
func (ch *ChannelHandle) Data() (*ChannelData, error) {
	return nil, nil
}

// Continue tells Asterisk to return the channel to the dialplan
func (ch *ChannelHandle) Continue(context, extension string, priority int) error {
	return nil
}

// Play initiates playback of the specified media uri
// to the channel, returning the Playback handle
//func (ch *ChannelHandle) Play(id string, mediaURI string) (ph *PlaybackHandle, err error) {
//	return nil, nil
//}

// Record records the channel to the given filename
//func (ch *ChannelHandle) Record(name string, opts *RecordingOptions) (*LiveRecordingHandle, error) {
//	return nil, nil//
//}

//---
// Hangup Operations
//---

// Busy hangs up the channel with the "busy" cause code
func (ch *ChannelHandle) Busy() error {
	return nil
}

// Congestion hangs up the channel with the congestion cause code
func (ch *ChannelHandle) Congestion() error {
	return nil
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
	updated, err := ch.Data()
	if err != nil {
		return false, err
	}

	return strings.ToLower(updated.State) == "up", nil
}

// Ring Operations
// --

// Ring indicates ringing to the channel
func (ch *ChannelHandle) Ring() error {
	return nil
}

// StopRing stops ringing on the channel
func (ch *ChannelHandle) StopRing() error {
	return nil
}

// Mute operations
// --

// Mute mutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Mute(dir arioptions.Direction) (err error) {
	if dir == "" {
		dir = arioptions.DirectionIn
	}

	return nil
}

// Unmute unmutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Unmute(dir arioptions.Direction) (err error) {
	if dir == "" {
		dir = arioptions.DirectionIn
	}

	return nil
}

// Hold operations
// --

// Hold puts the channel on hold
func (ch *ChannelHandle) Hold() error {
	return nil
}

// StopHold retrieves the channel from hold
func (ch *ChannelHandle) StopHold() error {
	return nil
}

// Music on hold operations
// --

// MOH plays music on hold of the given class
// to the channel
func (ch *ChannelHandle) MOH(mohClass string) error {
	return nil
}

// StopMOH stops playing of music on hold to the channel
func (ch *ChannelHandle) StopMOH() error {
	return nil
}

// ----

// GetVariable returns the value of a channel variable
func (ch *ChannelHandle) GetVariable(name string) (string, error) {
	return "", nil
}

// SetVariable sets the value of a channel variable
func (ch *ChannelHandle) SetVariable(name, value string) error {
	return nil
}

// --
// Misc
// --

// Originate creates (and dials) a new channel using the present channel as its Originator.
func (ch *ChannelHandle) Originate(req arioptions.OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return nil, nil
}

// StageOriginate stages an originate (channel creation and dial) to be Executed later.
func (ch *ChannelHandle) StageOriginate(req arioptions.OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return nil, nil
}

// Create creates (but does not dial) a new channel, using the present channel as its Originator.
func (ch *ChannelHandle) Create(req ChannelCreateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return nil, nil
}

// Dial dials a created channel.  `caller` is the optional
// channel ID of the calling party (if there is one).  Timeout
// is the length of time to wait before the dial is answered
// before aborting.
func (ch *ChannelHandle) Dial(caller string, timeout time.Duration) error {
	return nil
}

// Snoop spies on a specific channel, creating a new snooping channel placed into the given app
func (ch *ChannelHandle) Snoop(snoopID string, opts *SnoopOptions) (*ChannelHandle, error) {
	return nil, nil
}

// StageSnoop stages a `Snoop` operation
func (ch *ChannelHandle) StageSnoop(snoopID string, opts *SnoopOptions) (*ChannelHandle, error) {
	return nil, nil
}

// StageExternalMedia creates a new non-telephony external media channel,
// when `Exec`ed, by which audio may be sent or received.  The stage version
// of this command will not actually communicate with Asterisk until Exec is
// called on the returned ExternalMedia channel.
//func (ch *ChannelHandle) StageExternalMedia(opts ExternalMediaOptions) (*ChannelHandle, error) {/
//	return nil
//}

// ExternalMedia creates a new non-telephony external media channel by which audio may be sent or received
//func (ch *ChannelHandle) ExternalMedia(opts ExternalMediaOptions) (*ChannelHandle, error) {
//	return nil
//}

// Silence operations
// --

// Silence plays silence to the channel
func (ch *ChannelHandle) Silence() error {
	return nil
}

// StopSilence stops silence to the channel
func (ch *ChannelHandle) StopSilence() error {
	return nil
}

// DTMF
// --

// SendDTMF sends the DTMF information to the server
func (ch *ChannelHandle) SendDTMF(dtmf string, opts *arioptions.DTMFOptions) error {
	return nil
}

// UserEvent sends user-event to AMI channel subscribers
//func (ch *ChannelHandle) UserEvent(key *Key, ue *ChannelUserevent) error {
//	return nil
//}

type ChannelCreateRequest struct {
	// Endpoint is the target endpoint for the dial
	Endpoint string `json:"endpoint"`

	// App is the name of the Stasis application to execute on connection
	App string `json:"app"`

	// AppArgs is the set of (comma-separated) arguments for the Stasis App
	AppArgs string `json:"appArgs,omitempty"`

	// ChannelID is the ID to give to the newly-created channel
	ChannelID string `json:"channelId,omitempty"`

	// OtherChannelID is the ID of the second created channel (when creating Local channels)
	OtherChannelID string `json:"otherChannelId,omitempty"`

	// Originator is the unique ID of the calling channel, for which this new channel-dial is being created
	Originator string `json:"originator,omitempty"`

	// Formats is the comma-separated list of valid codecs to allow for the new channel, in the case that
	// the Originator is not specified
	Formats string `json:"formats,omitempty"`
}

// SnoopOptions enumerates the non-required arguments for the snoop operation
type SnoopOptions struct {
	// App is the ARI application into which the newly-created Snoop channel should be dropped.
	App string `json:"app"`

	// AppArgs is the set of arguments to pass with the newly-created Snoop channel's entry into ARI.
	AppArgs string `json:"appArgs,omitempty"`

	// Spy describes the direction of audio on which to spy (none, in, out, both).
	// The default is 'none'.
	Spy arioptions.Direction `json:"spy,omitempty"`

	// Whisper describes the direction of audio on which to send (none, in, out, both).
	// The default is 'none'.
	Whisper arioptions.Direction `json:"whisper,omitempty"`
}

// ExternalMediaOptions describes the parameters to the externalMedia channel creation operation
type ExternalMediaOptions struct {
	// ChannelID specifies the channel ID to be used for the external media channel.  This parameter is optional and if not specified, a randomly-generated channel ID will be used.
	ChannelID string `json:"channelId"`

	// App is the ARI Application to which the newly-created external media channel should be placed.  This parameter is optional and if not specified, the current application will be used.
	App string `json:"app"`

	// ExternalHost specifies the <host>:<port> of the external host to which the external media channel will be connected.  This parameter is MANDATORY and has no default.
	ExternalHost string `json:"external_host"`

	// Encapsulation specifies the payload encapsulation which should be used.  Options include:  'rtp'.  This parameter is optional and if not specified, 'rtp' will be used.
	Encapsulation string `json:"encapsulation"`

	// Transport specifies the connection type to be used to communicate to the external server.  Options include 'udp'.  This parameter is optional and if not specified, 'udp' will be used.
	Transport string `json:"transport"`

	// ConnectionType defined the directionality of the network connection.  Options include 'client' and 'server'.  This parameter is optional and if not specified, 'client' will be used.
	ConnectionType string `json:"connection_type"`

	// Format specifies the codec to be used for the audio.  Options include 'slin16', 'ulaw' (and likely other codecs supported by Asterisk).  This parameter is MANDATORY and has not default.
	Format string `json:"format"`

	// Direction specifies the directionality of the audio stream.  Options include 'both'.  This parameter is optional and if not specified, 'both' will be used.
	Direction string `json:"direction"`

	// Variables defines the set of channel variables which should be bound to this channel upon creation.  This parameter is optional.
	Variables map[string]string `json:"variables"`
}

//Code taken from ari-proxy. NOt verbatim, but with enough similarities to make this work .
//Key was a copy from ari-proxy
