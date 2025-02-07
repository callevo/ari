package arioptions

// ChannelCreateRequest describes how a channel should be created, when
// using the separate Create and Dial calls.
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
	AppArgs string `json:"app_args,omitempty"`

	// Spy describes the direction of audio on which to spy (none, in, out, both).
	// The default is 'none'.
	Spy Direction `json:"spy,omitempty"`

	// Whisper describes the direction of audio on which to send (none, in, out, both).
	// The default is 'none'.
	Whisper Direction `json:"whisper,omitempty"`
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
	Direction Direction `json:"direction"`

	// Variables defines the set of channel variables which should be bound to this channel upon creation.  This parameter is optional.
	Variables map[string]string `json:"variables"`
}
