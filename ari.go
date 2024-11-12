package ari

import (
	"time"

	commands "github.com/callevo/ari/command"
	"github.com/callevo/ari/logs"
	"github.com/callevo/ari/messagebus"
	"github.com/callevo/ari/proxy"
	"github.com/callevo/ari/requests"
	"github.com/lrita/cmap"
	nats "github.com/nats-io/nats.go"
)

type ARIClient struct {
	Application    string
	ConnectionName string
	NATSUrl        string

	announceSubs *nats.Subscription
	proxysubs    *nats.Subscription

	sbus         *messagebus.NatsBus
	_ast_cluster cmap.Cmap
}

func NewClient() *ARIClient {
	return &ARIClient{}
}

func NatsReconnect(nc *nats.Conn) error {
	return nil
}

func (a *ARIClient) Publish(cmd commands.Command) error {
	return a.sbus.PublishCommand(a.ConnectionName+".command", cmd)
}

func (a *ARIClient) Listen(opts *Options, natsURL string) error {
	logs.TLogger.Debug().Msg("Entering in listening mode")

	a.NATSUrl = natsURL
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

	a.sbus = messagebus.NewNatsBus(cfg)

	err := a.sbus.Connect()
	if err != nil {
		return err
	}

	logs.TLogger.Debug().Msg("subscribing to announce")
	a.announceSubs, err = a.sbus.SubscribeAnnounce(a.ConnectionName+".announce.*", func(o *proxy.Announcement) {
		logs.TLogger.Debug().Msgf("O: %+v", o)

		if o.Node != "" {
			// we need to look in the list
			requestTopic := a.ConnectionName + "." + a.Application + ".get." + o.Node
			astInfo := requests.NewAsteriskInfoRequest()
			astInfo.SetAsteriskID(o.Node)

			p, err := a.sbus.Request(requestTopic, astInfo)
			if err != nil {
				logs.TLogger.Debug().Msgf("error!! %+v", err)

				return
			}

			logs.TLogger.Debug().Msgf("O: %+v", p)
		}
	})
	if err != nil {
		logs.TLogger.Debug().Msgf("error!! %+v", err)

		return err
	}

	logs.TLogger.Debug().Msg("subscribing to events")
	a.proxysubs, err = a.sbus.SubscribeEvent(a.ConnectionName+"."+a.Application+".events.>", func(o *proxy.AriEvent) {
		logs.TLogger.Debug().Msgf("O: %+v", o)

		switch o.GetType() {

		}

	})
	if err != nil {
		logs.TLogger.Debug().Msgf("error!! %+v", err)

		return err
	}

	return nil
}

func (a *ARIClient) Close() {
	a.sbus.Close()
}

type Options struct {
	// Application is the the name of this ARI application
	Application string

	ConnectionName string
}
