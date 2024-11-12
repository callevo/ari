package proxy

type Events interface {
	GetType() string
	GetAsteriskID() string
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

type CallerInfo struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type DialplanInfo struct {
	Context  string `json:"context"`
	Exten    string `json:"exten"`
	Priority int    `json:"priority"`
	AppName  string `json:"app_name"`
	AppData  string `json:"app_data"`
}

type ChanneInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	Protocol     string `json:"protocol"`
	Caller       CallerInfo
	Connected    CallerInfo
	Accountcode  string `json:"accountcode"`
	Dialplan     DialplanInfo
	Creationtime string `json:"creationtime"`
	Language     string `json:"language"`
}

type AriEvent struct {
	Type        string   `json:"type"`
	AsteriskID  string   `json:"asteriskid"`
	Application string   `json:"application"`
	TimeStamp   string   `json:"timestamp"`
	Args        []string `json:"args"`
	Cause       int      `json:"cause,omitempty"`
	Channel     ChanneInfo
}

func (evt *AriEvent) GetType() string {
	return evt.Type
}

func (evt *AriEvent) GetApp() string {
	return evt.Application
}

func (evt *AriEvent) GetAsteriskID() string {
	return evt.AsteriskID
}
