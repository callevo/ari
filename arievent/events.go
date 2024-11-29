package arievent

import "github.com/callevo/ari/channel"

type Events interface {
	GetType() string
	GetNode() string

	// IsPropagationStopped informs weather the event should
	// be further propagated or not
	IsPropagationStopped() bool

	// StopPropagation makes the event no longer
	// propagate.
	StopPropagation()
}

type EventType string

type StasisEvent struct {
	Type            EventType `json:"type"`
	Node            string    `json:"asterisk_id"`
	Application     string    `json:"application"`
	TimeStamp       string    `json:"timestamp"`
	Args            []string  `json:"args"`
	Cause           int       `json:"cause,omitempty"`
	Channel         channel.ChannelData
	stopPropagation bool
}

func (evt *StasisEvent) GetType() EventType {
	return evt.Type
}

func (evt *StasisEvent) GetApp() string {
	return evt.Application
}

func (evt *StasisEvent) GetNode() string {
	return evt.Node
}

func (evt *StasisEvent) StopPropagation(p bool) {
	evt.stopPropagation = true
}

func (evt *StasisEvent) IsPropagationStopped() bool {
	return evt.stopPropagation
}

var (
	ApplicationMoveFailed    EventType = "ApplicationMoveFailed"
	ApplicationReplaced      EventType = "ApplicationReplaced"
	BridgeAttendedTransfer   EventType = "BridgeAttendedTransfer"
	BridgeBlindTransfer      EventType = "BridgeBlindTransfer"
	BridgeCreated            EventType = "BridgeCreated"
	BridgeDestroyed          EventType = "BridgeDestroyed"
	BridgeMerged             EventType = "BridgeMerged"
	BridgeVideoSourceChanged EventType = "BridgeVideoSourceChanged"
	ChannelCallerId          EventType = "ChannelCallerId"
	ChannelConnectedLine     EventType = "ChannelConnectedLine"
	ChannelCreated           EventType = "ChannelCreated"
	ChannelDestroyed         EventType = "ChannelDestroyed"
	ChannelDialplan          EventType = "ChannelDialplan"
	ChannelDtmfReceived      EventType = "ChannelDtmfReceived"
	ChannelEnteredBridge     EventType = "ChannelEnteredBridge"
	ChannelHangupRequest     EventType = "ChannelHangupRequest"
	ChannelHold              EventType = "ChannelHold"
	ChannelLeftBridge        EventType = "ChannelLeftBridge"
	ChannelStateChange       EventType = "ChannelStateChange"
	ChannelTalkingFinished   EventType = "ChannelTalkingFinished"
	ChannelTalkingStarted    EventType = "ChannelTalkingStarted"
	ChannelUnhold            EventType = "ChannelUnhold"
	ChannelUserevent         EventType = "ChannelUserevent"
	ChannelVarset            EventType = "ChannelVarset"
	ContactInfo              EventType = "ContactInfo"
	ContactStatusChange      EventType = "ContactStatusChange"
	DeviceStateChanged       EventType = "DeviceStateChanged"
	Dial                     EventType = "Dial"
	EndpointStateChange      EventType = "EndpointStateChange"
	Event                    EventType = "Event"
	Message                  EventType = "Message"
	MissingParams            EventType = "MissingParams"
	Peer                     EventType = "Peer"
	PeerStatusChange         EventType = "PeerStatusChange"
	PlaybackContinuing       EventType = "PlaybackContinuing"
	PlaybackFinished         EventType = "PlaybackFinished"
	PlaybackStarted          EventType = "PlaybackStarted"
	RecordingFailed          EventType = "RecordingFailed"
	RecordingFinished        EventType = "RecordingFinished"
	RecordingStarted         EventType = "RecordingStarted"
	StasisEnd                EventType = "StasisEnd"
	StasisStart              EventType = "StasisStar1t"
	TextMessageReceived      EventType = "TextMessageReceived"
)
