package requests

type Request interface {
	GetName() string
	GetAsteriskID() string
	SetAsteriskID(id string)
}

type AsteriskInfoRequest struct {
	Name       string `json:"name"`
	ASteriskID string `json:"asteriskid,omitempty"`
}

func NewAsteriskInfoRequest() *AsteriskInfoRequest {
	r := &AsteriskInfoRequest{
		Name: "asteriskinfo",
	}

	return r
}

func (a *AsteriskInfoRequest) GetName() string {
	return a.Name
}

func (a *AsteriskInfoRequest) GetAsteriskID() string {
	return a.ASteriskID
}

func (a *AsteriskInfoRequest) SetAsteriskID(id stgring) {
	a.ASteriskID = id
}
