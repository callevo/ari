package proxy

type Events interface {
}

type EventType int

const (
	StasisStart EventType = iota + 1 // EnumIndex = 1
	StasisEnd
	ChannelHangupRequest
)

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (w EventType) EnumIndex() int {
	return int(w)
}
