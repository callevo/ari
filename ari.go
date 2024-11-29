package ari

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/callevo/ari/arievent"
	"github.com/callevo/ari/asterisk"
	"github.com/callevo/ari/bridge"
	"github.com/callevo/ari/channel"
	"github.com/callevo/ari/cluster"
	"github.com/callevo/ari/dispatcher"
	"github.com/callevo/ari/key"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/messagebus"
	"github.com/callevo/ari/requests"
	"github.com/callevo/ari/response"
	"github.com/lrita/cmap"
	nats "github.com/nats-io/nats.go"
	"github.com/rotisserie/eris"
)

// ErrNil indicates that the request returned an empty response
var ErrNil = eris.New("Nil")

type ARIClient struct {
	Application    string
	ConnectionName string
	NATSUrl        string

	announceSubs *nats.Subscription
	proxysubs    *nats.Subscription

	sbus              *messagebus.NatsBus
	_ast_cluster      cmap.Cmap
	_dynSubscriptions cmap.Cmap

	_dispatcher *dispatcher.EventDispatcher

	mu sync.Mutex
}

func NewClient() *ARIClient {
	return &ARIClient{}
}

func NatsReconnect(nc *nats.Conn) error {
	return nil
}

// Subject returns the communication subject for the given parameters
func Subject(prefix, appName, class, asterisk string) (ret string) {
	ret = prefix + "."
	if appName != "" {
		ret += (appName + "." + class)
		if asterisk != "" {
			ret += "." + asterisk
		}
	}
	return
}

type StasisHandler func(*ARIClient, *channel.ChannelHandle, *arievent.StasisEvent)

func (a *ARIClient) Listen(ctx context.Context, opts *Options, exechandler StasisHandler) error {
	logs.TLogger.Debug().Msg("Entering in listening mode")

	a.NATSUrl = opts.NatsUrl
	a.ConnectionName = opts.ConnectionName
	a.Application = opts.Application

	cfg := messagebus.Config{
		URL:            a.NATSUrl,
		NatsTimeout:    10 * time.Second,
		RequestTimeout: 3 * time.Second,
		ConnectionName: a.ConnectionName,
		PingInterval:   20 * time.Second,
		MaxReconnects:  10,
		MaxPing:        3,
	}

	a._dispatcher = dispatcher.NewDispatcher()

	a.sbus = messagebus.NewNatsBus(cfg)

	err := a.sbus.Connect()
	if err != nil {
		return err
	}

	logs.TLogger.Debug().Msg("subscribing to announce")
	a.announceSubs, err = a.sbus.SubscribeAnnounce(a.ConnectionName+".announce.*", func(o *cluster.Announcement) {
		//logs.TLogger.Debug().Msgf("O: %+v", o)

		if o.Node != "" {
			// we need to look in the list
			/*
				requestTopic := a.ConnectionName + "." + a.Application + ".get." + o.Node
				astInfo := requests.NewAsteriskInfoRequest()
				astInfo.SetAsteriskID(o.Node)

				p, err := a.sbus.Request(requestTopic, astInfo)
				if err != nil {
					logs.TLogger.Debug().Msgf("error!! %+v", err)

					return
				}

				var rsp map[string]interface{}
				err = json.Unmarshal(p.Data, &rsp)
				if err != nil {
					logs.TLogger.Debug().Msgf("error!! %+v", err)

					return
				}

				logs.TLogger.Debug().Msgf("Rsp: %+v", rsp)
			*/
		}
	})
	if err != nil {
		logs.TLogger.Debug().Msgf("error!! %+v", err)

		return err
	}

	logs.TLogger.Debug().Msg("Queue subscribing to stasisstart events %s")
	a.proxysubs, err = a.sbus.SubscribeEvent(a.ConnectionName+"."+a.Application+".*.*.stasisstart.>", func(o *arievent.StasisEvent) {
		logs.TLogger.Debug().Msgf("O: %+v", o)

		// We need to dispatch Event

		a._dispatcher.Dispatch(o)

		k := key.NewKey(key.ChannelKey, o.Channel.GetID(), key.WithApp(o.Application), key.WithNode(o.Node))

		h := channel.NewChannelHandle(k, &ichannel{c: a}, nil)

		go exechandler(a, h, o)

		channelTopic := a.ConnectionName + "." + a.Application + "." + o.Node + "." + strings.ReplaceAll(o.Channel.ID, ".", "#")
		logs.TLogger.Debug().Msgf("subscribing client to %s", channelTopic)
		dynSub, err := a.sbus.DynSubscription(channelTopic, func(o *arievent.StasisEvent) {

			logs.TLogger.Debug().Msgf("O: %+v", o)

			//dispatching the event to the listeners
			a._dispatcher.Dispatch(o)

			switch o.GetType() {
			//case arievent.ApplicationMoveFailed:
			//case arievent.ApplicationReplaced:
			//case arievent.BridgeAttendedTransfer:
			//case arievent.BridgeBlindTransfer:
			//case arievent.BridgeCreated:
			//case arievent.BridgeDestroyed:
			//case arievent.BridgeMerged:
			//case arievent.BridgeVideoSourceChanged:
			//case arievent.ChannelCallerId:
			//case arievent.ChannelConnectedLine:
			//case arievent.ChannelCreated:
			//case arievent.ChannelDestroyed:
			//case arievent.ChannelDialplan:
			//case arievent.ChannelDtmfReceived:
			//case arievent.ChannelEnteredBridge:
			//case arievent.ChannelHangupRequest:
			//case arievent.ChannelHold:
			//case arievent.ChannelLeftBridge:
			//case arievent.ChannelStateChange:
			//case arievent.ChannelTalkingFinished:
			//case arievent.ChannelTalkingStarted:
			//case arievent.ChannelUnhold:
			//case arievent.ChannelUserevent:
			//case arievent.ChannelVarset:
			//case arievent.ContactInfo:
			//case arievent.ContactStatusChange:
			//case arievent.DeviceStateChanged:
			//case arievent.Dial:
			//case arievent.EndpointStateChange:
			//case arievent.Event:
			//case arievent.Message:
			//case arievent.MissingParams:
			//case arievent.PeerStatusChange:
			//case arievent.PlaybackContinuing:
			//case arievent.PlaybackFinished:
			//case arievent.PlaybackStarted:
			//case arievent.RecordingFailed:
			//case arievent.RecordingFinished:
			//case arievent.RecordingStarted:
			//case arievent.StasisStart:
			//case arievent.TextMessageReceived:
			case arievent.StasisEnd:

				channelTopic := a.ConnectionName + "." + a.Application + "." + o.Node + "." + strings.ReplaceAll(o.Channel.ID, ".", "#")
				if myDynSub, ok := a._dynSubscriptions.Load(channelTopic); ok {
					logs.TLogger.Debug().Msgf("call finished we need to drain and unscrubscribe")
					myDynSub.(*nats.Subscription).Drain()

					a._dynSubscriptions.Delete(channelTopic)
				}
			default:

			}
		})
		if err != nil {
			logs.TLogger.Debug().Msgf("error!! %+v", err)

			return
		}

		a._dynSubscriptions.Store(channelTopic, dynSub)
	})
	if err != nil {
		logs.TLogger.Debug().Msgf("error!! %+v", err)

		return eris.Wrap(err, "error creating dynamic subscription for topic")
	}

	<-ctx.Done()

	return ctx.Err()
}

func (a *ARIClient) Channel() channel.Channel {
	return &ichannel{a}
}

func (a *ARIClient) Asterisk() asterisk.Asterisk {
	return &iasterisk{a}
}

func (a *ARIClient) Bridge() bridge.Bridge {
	return &ibridge{a}
}

func (a *ARIClient) Dispatcher() *dispatcher.EventDispatcher {
	return a._dispatcher
}

func (a *ARIClient) Close() {
	a.sbus.Close()
}

type Options struct {
	// Application is the the name of this ARI application
	Application string

	ConnectionName string

	NatsUrl string
}

func (c *ARIClient) commandRequest(req *requests.Request) error {
	resp, err := c.makeRequest("command", req)
	if err != nil {
		return err
	}
	return resp.Err()
}

func (c *ARIClient) completeCoordinates(req *requests.Request) bool {
	if req == nil || req.Key == nil {
		return false
	}

	// coordinates are complete if we have both app and node
	return req.Key.App != "" && req.Key.Node != ""
}

func (c *ARIClient) makeRequest(class string, req *requests.Request) (*response.Response, error) {
	//var resp response.Response
	var err error

	if !c.completeCoordinates(req) {
		return nil, err
	}

	logs.TLogger.Debug().Msgf("Sending request to %s", c.subject(class, req))
	return c.sbus.Request(c.subject(class, req), req)
}

func (c *ARIClient) subject(class string, req *requests.Request) string {
	if req == nil || req.Key == nil {
		return Subject(c.ConnectionName, c.Application, class, "")
	}
	return Subject(c.ConnectionName, req.Key.App, class, req.Key.Node)
}

func (c *ARIClient) getRequest(req *requests.Request) (*key.Key, error) {
	resp, err := c.makeRequest("get", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Key == nil {
		return nil, ErrNil
	}
	return resp.Key, nil
}

func (c *ARIClient) dataRequest(req *requests.Request) (*response.EntityData, error) {
	resp, err := c.makeRequest("data", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Data == nil {
		return nil, ErrNil
	}
	return resp.Data, nil
}

func (c *ARIClient) createRequest(req *requests.Request) (*key.Key, error) {
	resp, err := c.makeRequest("create", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Key == nil {
		return nil, ErrNil
	}
	return resp.Key, nil
}
