package messagebus

import (
	"encoding/json"
	"fmt"
	"time"

	arievent "github.com/callevo/ari/arievent"
	cluster "github.com/callevo/ari/cluster"
	logs "github.com/callevo/ari/logs"
	requests "github.com/callevo/ari/requests"
	"github.com/callevo/ari/response"
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

// PublishAnnounce sends announce message
func (n *NatsBus) PublishAnnounce(topic string, msg *cluster.Announcement) error {
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
type AnnounceHandler func(o *cluster.Announcement)

// EventHandler Handles Events
type EventHandler func(o *arievent.StasisEvent)

// SubscribeAnnounce subscribe announce messages
func (n *NatsBus) SubscribeAnnounce(topic string, callback AnnounceHandler) (*nats.Subscription, error) {
	logs.TLogger.Debug().Msgf("Subscribing to %s", topic)
	return n.conn.Subscribe(topic, func(msg *nats.Msg) {
		evt := cluster.Announcement{}

		logs.TLogger.Debug().Msgf("We got %s", msg.Data)
		err := json.Unmarshal(msg.Data, &evt)
		if err != nil {
			return
		}

		callback(&evt)
	})
}

// ListenQueue is the queue group to use for distributing arieventStart events to Listeners.
var ListenQueue = "AsteriskARIProxyDistributionQueue"

func (n *NatsBus) SubscribeEvent(topic string, callback EventHandler) (*nats.Subscription, error) {
	logs.TLogger.Debug().Msgf("Subscribing to %s", topic)

	return n.conn.QueueSubscribe(topic, ListenQueue, func(msg *nats.Msg) {

		evt := arievent.StasisEvent{}

		logs.TLogger.Debug().Msgf("We got %s", (string)(msg.Data))

		err := json.Unmarshal(msg.Data, &evt)
		if err != nil {
			return
		}

		if callback != nil {
			callback(&evt)
		}
	})
}

func (n *NatsBus) DynSubscription(topic string, callback EventHandler) (*nats.Subscription, error) {
	logs.TLogger.Debug().Msgf("Subscribing to %s", topic+".>")

	return n.conn.Subscribe(topic+".>", func(msg *nats.Msg) {
		evt := arievent.StasisEvent{}

		err := json.Unmarshal(msg.Data, &evt)
		if err != nil {
			return
		}

		if callback != nil {
			callback(&evt)
		}
	})
}

func (n *NatsBus) Request(topic string, r *requests.Request) (*response.Response, error) {

	if n.conn == nil {
		return nil, fmt.Errorf("nil connection")
	}

	b, err := json.Marshal(r)
	if err != nil {
		logs.TLogger.Debug().Msgf("err %s", err)

		return nil, err
	}

	msg, err := n.conn.Request(topic, b, 3*time.Second)
	if err != nil {
		logs.TLogger.Debug().Msgf("err %s", err)

		return nil, err
	}

	resp := &response.Response{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		logs.TLogger.Debug().Msgf("err %s", err)

		return nil, err
	}

	logs.TLogger.Debug().Msgf("we got this response: %+v", resp)

	return resp, nil
}
