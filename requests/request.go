package requests

import (
	"time"

	"github.com/callevo/ari/arioptions"
	"github.com/callevo/ari/key"
)

// Request describes a request which is sent from an ARI proxy Client to an ARI proxy Server
type Request struct {
	// Kind indicates the type of request
	Kind string `json:"kind"`

	// Key is the key or key filter on which this request should be processed
	Key *key.Key `json:"key"`

	AsteriskConfig         *AsteriskConfig         `json:"asterisk_config,omitempty"`
	AsteriskLoggingChannel *AsteriskLoggingChannel `json:"asterisk_logging_channel,omitempty"`
	AsteriskVariableSet    *AsteriskVariableSet    `json:"asterisk_variable_set,omitempty"`

	BridgeAddChannel    *BridgeAddChannel    `json:"bridge_add_channel,omitempty"`
	BridgeCreate        *BridgeCreate        `json:"bridge_create,omitempty"`
	BridgeMOH           *BridgeMOH           `json:"bridge_moh,omitempty"`
	BridgePlay          *BridgePlay          `json:"bridge_play,omitempty"`
	BridgeRecord        *BridgeRecord        `json:"bridge_record,omitempty"`
	BridgeRemoveChannel *BridgeRemoveChannel `json:"bridge_remove_channel,omitempty"`
	BridgeVideoSource   *BridgeVideoSource   `json:"bridge_video_source,omitempty"`

	ChannelCreate        *ChannelCreate        `json:"channel_create,omitempty"`
	ChannelContinue      *ChannelContinue      `json:"channel_continue,omitempty"`
	ChannelDial          *ChannelDial          `json:"channel_dial,omitempty"`
	ChannelHangup        *ChannelHangup        `json:"channel_hangup,omitempty"`
	ChannelMOH           *ChannelMOH           `json:"channel_moh,omitempty"`
	ChannelMute          *ChannelMute          `json:"channel_mute,omitempty"`
	ChannelOriginate     *ChannelOriginate     `json:"channel_originate,omitempty"`
	ChannelPlay          *ChannelPlay          `json:"channel_play,omitempty"`
	ChannelRecord        *ChannelRecord        `json:"channel_record,omitempty"`
	ChannelSendDTMF      *ChannelSendDTMF      `json:"channel_send_dtmf,omitempty"`
	ChannelSnoop         *ChannelSnoop         `json:"channel_snoop,omitempty"`
	ChannelExternalMedia *ChannelExternalMedia `json:"channel_external_media,omitempty"`
	ChannelVariable      *ChannelVariable      `json:"channel_variable,omitempty"`
	//ChannelUserevent     *ChannelUserevent     `json:"channel_user_event,omitempty"`

	DeviceStateUpdate *DeviceStateUpdate `json:"device_state_update,omitempty"`

	EndpointListByTech *EndpointListByTech `json:"endpoint_list_by_tech,omitempty"`

	MailboxUpdate *MailboxUpdate `json:"mailbox_update,omitempty"`

	PlaybackControl *PlaybackControl `json:"playback_control,omitempty"`

	RecordingStoredCopy *RecordingStoredCopy `json:"recording_stored_copy,omitempty"`

	SoundList *SoundList `json:"sound_list,omitempty"`
}

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

// ApplicationSubscribe describes a request to subscribe/unsubscribe a particular ARI application to an EventSource
type ApplicationSubscribe struct {
	// EventSource is the ARI event source to which the subscription is requested.  This should be one of:
	//  - channel:<channelId>
	//  - bridge:<bridgeId>
	//  - endpoint:<tech>/<resource> (e.g. SIP/102)
	//  - deviceState:<deviceName>
	EventSource string `json:"event_source"`
}

// AsteriskVariableSet is the request type for setting an asterisk variable
type AsteriskVariableSet struct {
	// Value is the value to set
	Value string `json:"value"`
}

// BridgeAddChannel is the request type for adding a channel to a bridge
type BridgeAddChannel struct {
	// Channel is the channel ID to add to the bridge
	Channel string `json:"channel"`

	// AbsorbDTMF indicates that DTMF coming from this channel will not be passed through to the bridge
	AbsorbDTMF bool `json:"absorbDTMF,omitempty"`

	// Mute indicates that the channel should be muted, preventing audio from it passing through to the bridge
	Mute bool `json:"mute,omitempty"`

	// Role indicates the channel's role in the bridge
	Role string `json:"role,omitempty"`
}

// BridgeCreate is the request type for creating a bridge
type BridgeCreate struct {
	// Type is the comma-separated list of bridge type attributes (mixing,
	// holding, dtmf_events, proxy_media).  If not set, the default (mixing)
	// will be used.
	Type string `json:"type"`

	// Name is the name to assign to the bridge (optional)
	Name string `json:"name,omitempty"`
}

// BridgeMOH is the request type for playing Music on Hold to a bridge
type BridgeMOH struct {
	// Class is the Music On Hold class to be played
	Class string `json:"class"`
}

// BridgePlay is the request type for playing audio on the bridge
type BridgePlay struct {
	// PlaybackID is the unique identifier for this playback
	PlaybackID string `json:"playback_id"`

	// MediaURI is the URI from which to obtain the playback media
	MediaURI string `json:"media_uri"`
}

// BridgeRecord is the request for recording a bridge
type BridgeRecord struct {
	//	// Name is the name for the recording
	Name string `json:"name"`

	// Options is the list of recording Options
	Options *arioptions.RecordingOptions `json:"options,omitempty"`
}

// BridgeRemoveChannel is the request for removing a channel on the bridge
type BridgeRemoveChannel struct {
	// Channel is the name of the channel to remove
	Channel string `json:"channel"`
}

// BridgeVideoSource describes the details of a request to set the video source of a bridge explicitly
type BridgeVideoSource struct {
	// Channel is the name of the channel to use as the explicit video source
	Channel string `json:"channel"`
}

// ChannelCreate describes a request to create a new channel
type ChannelCreate struct {
	// ChannelCreateRequest is the request for creating the channel
	ChannelCreateRequest ChannelCreateRequest `json:"channel_create_request"`
}

// ChannelContinue describes a request to continue an ARI application
type ChannelContinue struct {
	// Context is the context into which the channel should be continued
	Context string `json:"context"`

	// Extension is the extension into which the channel should be continued
	Extension string `json:"extension"`

	// Priority is the priority at which the channel should be continued
	Priority int `json:"priority"`
}

// ChannelDial describes a request to dial
type ChannelDial struct {
	// Caller is the channel ID of the "caller" channel; if specified, the media parameters of the dialing channel will be matched to the "caller" channel.
	Caller string `json:"caller"`

	// Timeout is the maximum time which should be allowed for the dial to complete
	Timeout time.Duration `json:"timeout"`
}

// ChannelHangup is the request for hanging up a channel
type ChannelHangup struct {
	// Reason is the reason the channel is being hung up
	Reason string `json:"reason"`
}

// ChannelMOH is the request playing hold on music on a channel
type ChannelMOH struct {
	// Music is the music to play
	Music string `json:"music"`
}

// ChannelMute is the request for muting or unmuting a channel
type ChannelMute struct {
	// Direction is the direction to mute
	Direction arioptions.Direction `json:"direction,omitempty"`
}

// ChannelOriginate is the request for creating a channel
type ChannelOriginate struct {
	// OriginateRequest contains the information for originating a channel
	OriginateRequest OriginateRequest `json:"originate_request"`
}

// ChannelPlay is the request for playing audio on a channel
type ChannelPlay struct {
	// PlaybackID is the unique identifier for this playback
	PlaybackID string `json:"playback_id"`

	// MediaURI is the URI from which to obtain the playback media
	MediaURI string `json:"media_uri"`
}

// ChannelRecord is the request for recording a channel
type ChannelRecord struct {
	// Name is the name for the recording/
	Name string `json:"name"`

	// Options is the list of recording Options
	Options *arioptions.RecordingOptions `json:"options,omitempty"`
}

// ChannelSendDTMF is the request for sending a DTMF event to a channel
type ChannelSendDTMF struct {
	// DTMF is the series of DTMF inputs to send
	DTMF string `json:"dtmf"`

	// Options are the DTMF options
	Options *arioptions.DTMFOptions `json:"options,omitempty"`
}

// ChannelSnoop is the request for snooping on a channel
type ChannelSnoop struct {
	// SnoopID is the ID to use for the snoop channel which will be created.
	SnoopID string `json:"snoop_id"`

	// Options describe the parameters for the snoop session
	Options *arioptions.SnoopOptions `json:"options,omitempty"`
}

// ChannelExternalMedia describes the request for an external media channel
type ChannelExternalMedia struct {
	Options arioptions.ExternalMediaOptions `json:"options"`
}

// ChannelVariable is the request type to read or modify a channel variable
type ChannelVariable struct {
	// Name is the name of the channel variable
	Name string `json:"name"`

	// Value is the value to set to the channel variable
	Value string `json:"value,omitempty"`
}

// DeviceStateUpdate describes the request for updating the device state
type DeviceStateUpdate struct {
	// State is the new state of the device to set
	State string `json:"state"`
}

// EndpointListByTech describes the request for listing endpoints by technology
type EndpointListByTech struct {
	// Tech is the technology for the endpoint
	Tech string `json:"tech"`
}

// MailboxUpdate describes the request for updating a mailbox
type MailboxUpdate struct {
	// New is the number of New (unread) messages in the mailbox
	New int `json:"new"`

	// Old is the number of Old (read) messages in the mailbox
	Old int `json:"old"`
}

// PlaybackControl describes the request for performing a playback command
type PlaybackControl struct {
	// Command is the playback control command to run
	Command string `json:"command"`
}

// RecordingStoredCopy describes the request for copying a stored recording
type RecordingStoredCopy struct {
	// Destination is the destination location to copy to
	Destination string `json:"destination"`
}

// SoundList describes the request for listing the sounds
type SoundList struct {
	// Filters are the filters to apply when listing the sounds
	Filters map[string]string `json:"filters"`
}

// AsteriskConfig describes the request relating to asterisk configuration
type AsteriskConfig struct {
	// Tuples is the list of configuration tuples to update
	Tuples []arioptions.ConfigTuple `json:"tuples,omitempty"`
}

// AsteriskLoggingChannel describes a request relating to an asterisk logging channel
type AsteriskLoggingChannel struct {
	// Levels is the set of logging levels for this logging channel (comma-separated string)
	Levels string `json:"config"`
}

// ChannelUserevent - "User-generated event with additional user-defined fields in the object."
//type ChannelUserevent struct {
//	UserEvent arioptions.ChannelUserevent
//}
