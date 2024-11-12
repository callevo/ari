package messagebus

import (
	"encoding/json"
	"fmt"
	"time"

	commands "github.com/callevo/ari/command"
	logs "github.com/callevo/ari/logs"
	proxy "github.com/callevo/ari/proxy"
	requests "github.com/callevo/ari/requests"
	nats "github.com/nats-io/nats.go"
)

// DefaultReconnectionAttemts is the default number of reconnection attempts
// It implements a hard coded fault tolerance for a starting NATS cluster
const DefaultReconnectionAttemts = 5

// DefaultReconnectionWait is the default wating time between each reconnection
// attempt
const DefaultReconnectionWait = 5 * time.Second

// Config has general configuration for MessageBus
type Config struct {
	URL            string
	TimeoutRetries int
	NatsTimeout    time.Duration
	RequestTimeout time.Duration
	ConnectionName string
	PingInterval   time.Duration
	MaxReconnects  int
	MaxPing        int
}

type NatsBus struct {
	Config          Config
	conn            *nats.Conn
	ConnectedServer string
	ReconnHandler   nats.ConnHandler
}

// OptionNatsFunc options for RabbitMQ
type OptionNatsFunc func(n *NatsBus)

// NewNatsBus creates a NatsBus
func NewNatsBus(config Config, options ...OptionNatsFunc) *NatsBus {

	mbus := NatsBus{
		Config: config,
	}

	mbus.Config.MaxPing = 3
	mbus.Config.MaxReconnects = DefaultReconnectionAttemts
	mbus.Config.NatsTimeout = DefaultReconnectionWait * time.Second
	mbus.Config.PingInterval = 20 * time.Second

	for _, optfn := range options {
		optfn(&mbus)
	}

	return &mbus
}

// WithNatsConn binds an existing NATS connection
func WithNatsConn(nconn *nats.Conn) OptionNatsFunc {
	return func(n *NatsBus) {
		n.conn = nconn
	}
}

func WithNatsName(name string) OptionNatsFunc {
	return func(n *NatsBus) {
		n.Config.ConnectionName = name
	}
}

func WithReconnectionHandler(cb nats.ConnHandler) OptionNatsFunc {
	return func(n *NatsBus) {
		n.ReconnHandler = cb
	}
}

// Connect creates a NATS connection
func (n *NatsBus) Connect() error {
	var err error

	n.conn, err = nats.Connect(n.Config.URL,
		nats.Name(n.Config.ConnectionName),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			logs.TLogger.Debug().Msgf("Known servers: %v", nc.Servers())
			logs.TLogger.Debug().Msgf("Discovered servers: %v", nc.DiscoveredServers())
		}),
		nats.ReconnectWait(n.Config.NatsTimeout),
		nats.MaxReconnects(n.Config.MaxReconnects),
		nats.PingInterval(n.Config.PingInterval),
		nats.MaxPingsOutstanding(n.Config.MaxPing),
		//nats.NoEcho(),
		nats.DisconnectHandler(func(c *nats.Conn) {
			logs.TLogger.Debug().Msgf("Disconnected FROM %s", n.ConnectedServer)
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			if n.ReconnHandler != nil {
				n.ReconnHandler(c)
			}
		}),
	)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return err
	}

	return nil
}

func (n *NatsBus) PublishCommand(topic string, cmd commands.Command) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return n.conn.Publish(topic, b)
}

// PublishAnnounce sends announce message
func (n *NatsBus) PublishAnnounce(topic string, msg *proxy.Announcement) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(topic, b)
}

// Close closes the connection
func (n *NatsBus) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

// AnnounceHandler handles announce messages
type AnnounceHandler func(o *proxy.Announcement)

// EventHandler Handles Events
type EventHandler func(o *proxy.AriEvent)

// SubscribeAnnounce subscribe announce messages
func (n *NatsBus) SubscribeAnnounce(topic string, callback AnnounceHandler) (*nats.Subscription, error) {
	logs.TLogger.Debug().Msgf("Subscribing to %s", topic)
	return n.conn.Subscribe(topic, func(msg *nats.Msg) {
		evt := proxy.Announcement{}

		logs.TLogger.Debug().Msgf("We got %s", msg.Data)
		err := json.Unmarshal(msg.Data, &evt)
		if err != nil {
			return
		}

		callback(&evt)
	})
}

// ListenQueue is the queue group to use for distributing StasisStart events to Listeners.
var ListenQueue = "AsteriskARIProxyDistributionQueue"

func (n *NatsBus) SubscribeEvent(topic string, callback EventHandler) (*nats.Subscription, error) {
	logs.TLogger.Debug().Msgf("Subscribing to %s", topic)

	return n.conn.QueueSubscribe(topic, ListenQueue, func(msg *nats.Msg) {

		evt := proxy.AriEvent{}

		err := json.Unmarshal(msg.Data, &evt)
		if err != nil {
			return
		}

		if callback != nil {
			callback(&evt)
		}
	})
}

func (n *NatsBus) Request(topic string, r requests.Request) (*nats.Msg, error) {

	if n.conn == nil {
		return nil, fmt.Errorf("nil connection")
	}

	b, err := json.Marshal(r)
	if err != nil {
		logs.TLogger.Debug().Msgf("err %s", err)

		return nil, err
	}

	return n.conn.Request(topic, b, 3*time.Second)
}
