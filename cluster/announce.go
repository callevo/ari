package cluster

// Announcement describes the structure of an ARI proxy's announcement of availability on the network.  These are sent periodically and upon request (by a Ping).
type Announcement struct {
	// EventName
	EventName string `json:"name"`

	// Node indicates the Asterisk ID to which the proxy is connected
	Node string `json:"asteriskid"`

	// Application indicates the ARI application as which the proxy is connected
	Application string `json:"application"`
}
