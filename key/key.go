package key

const (
	// ApplicationKey is the key kind for ARI Application resources.
	ApplicationKey = "application"

	// BridgeKey is the key kind for the ARI Bridge resources.
	BridgeKey = "bridge"

	// ChannelKey is the key kind for the ARI Channel resource
	ChannelKey = "channel"

	// DeviceStateKey is the key kind for the ARI DeviceState resource
	DeviceStateKey = "devicestate"

	// EndpointKey is the key kind for the ARI Endpoint resource
	EndpointKey = "endpoint"

	// LiveRecordingKey is the key kind for the ARI LiveRecording resource
	LiveRecordingKey = "liverecording"

	// LoggingKey is the key kind for the ARI Logging resource
	LoggingKey = "logging"

	// MailboxKey is the key kind for the ARI Mailbox resource
	MailboxKey = "mailbox"

	// ModuleKey is the key kind for the ARI Module resource
	ModuleKey = "module"

	// PlaybackKey is the key kind for the ARI Playback resource
	PlaybackKey = "playback"

	// SoundKey is the key kind for the ARI Sound resource
	SoundKey = "sound"

	// StoredRecordingKey is the key kind for the ARI StoredRecording resource
	StoredRecordingKey = "storedrecording"

	// VariableKey is the key kind for the ARI Asterisk Variable resource
	VariableKey = "variable"
)

// KeyOptionFunc is a functional argument alias for providing options for ARI keys
type KeyOptionFunc func(Key) Key

type Key struct {
	Kind string `json:"kind,omitempty"`
	ID   string `json:"id,omitempty"`
	Node string `json:"node,omitempty"`
	App  string `json:"app,omitempty"`
}

// WithNode sets the given node identifier on the key.
func WithNode(node string) KeyOptionFunc {
	return func(key Key) Key {
		key.Node = node
		return key
	}
}

// WithApp sets the given node identifier on the key.
func WithApp(app string) KeyOptionFunc {
	return func(key Key) Key {
		key.App = app
		return key
	}
}

// NewKey builds a new key given the kind, identifier, and any optional arguments.
func NewKey(kind string, id string, opts ...KeyOptionFunc) *Key {
	k := Key{
		Kind: kind,
		ID:   id,
	}
	for _, o := range opts {
		k = o(k)
	}

	return &k
}

// NodeKey returns a key that is bound to the given application and node
func NodeKey(app, node string) *Key {
	return NewKey("", "", WithApp(app), WithNode(node))
}

// New returns a new key with the location information from the source key.
// This includes the App, the Node, and the Dialog.  the `kind` and `id`
// parameters are optional.  If kind is empty, the resulting key will not be
// typed.  If id is empty, the key will not be unique.
func (k *Key) New(kind, id string) *Key {
	n := NodeKey(k.App, k.Node)
	//n.Dialog = k.Dialog
	n.Kind = kind
	n.ID = id

	return n
}

func (m *Key) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *Key) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Key) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *Key) GetApp() string {
	if m != nil {
		return m.App
	}
	return ""
}
