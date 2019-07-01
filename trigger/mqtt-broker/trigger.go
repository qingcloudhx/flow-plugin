package mqtt_broker

import (
	"context"
	"strings"
	"time"

	"github.com/qingcloudhx/core/support/ssl"
	"github.com/qingcloudhx/gomqtt/broker"
	"github.com/qingcloudhx/gomqtt/packet"
	"github.com/qingcloudhx/gomqtt/transport"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Trigger is simple MQTT trigger
type Trigger struct {
	handlers  map[string]*clientHandler
	settings  *Settings
	logger    log.Logger
	options   *mqtt.ClientOptions
	client    mqtt.Client
	server    broker.Engine
	deviceCon *DeviceCon
}
type clientHandler struct {
	//client mqtt.Client
	handler  trigger.Handler
	settings *HandlerSettings
}
type Factory struct {
}

func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}

	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s}, nil
}

// Initialize implements trigger.Initializable.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.logger = ctx.Logger()

	settings := t.settings
	options := initClientOption(settings)
	t.options = options

	//client
	if strings.HasPrefix(settings.Broker, "ssl") {

		cfg := &ssl.Config{}

		if len(settings.SSLConfig) != 0 {
			err := cfg.FromMap(settings.SSLConfig)
			if err != nil {
				return err
			}

			if _, set := settings.SSLConfig["skipVerify"]; !set {
				cfg.SkipVerify = true
			}
			if _, set := settings.SSLConfig["useSystemCert"]; !set {
				cfg.UseSystemCert = true
			}
		} else {
			//using ssl but not configured, use defaults
			cfg.SkipVerify = true
			cfg.UseSystemCert = true
		}

		tlsConfig, err := ssl.NewClientTLSConfig(cfg)
		if err != nil {
			return err
		}

		options.SetTLSConfig(tlsConfig)
	}

	options.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		t.logger.Warnf("Recieved message on unhandled topic: %s", msg.Topic())
	})

	t.logger.Debugf("Client options: %v", options)

	t.handlers = make(map[string]*clientHandler)

	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		t.handlers[s.Topic] = &clientHandler{handler: handler, settings: s}
	}

	//server
	t.deviceCon = NewDeviceCon()
	return nil
}

func initClientOption(settings *Settings) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(settings.Broker)
	opts.SetClientID(settings.Id)
	opts.SetUsername(settings.Username)
	opts.SetPassword(settings.Password)
	opts.SetCleanSession(settings.CleanSession)
	opts.SetAutoReconnect(settings.AutoReconnect)

	if settings.Store != ":memory:" && settings.Store != "" {
		opts.SetStore(mqtt.NewFileStore(settings.Store))
	}

	if settings.KeepAlive != 0 {
		opts.SetKeepAlive(time.Duration(settings.KeepAlive) * time.Second)
	} else {
		opts.SetKeepAlive(2 * time.Second)
	}

	return opts
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() (err error) {

	client := mqtt.NewClient(t.options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	t.client = client

	for _, handler := range t.handlers {

		if token := client.Subscribe(handler.settings.Topic, byte(handler.settings.Qos), t.getHanlder(handler)); token.Wait() && token.Error() != nil {
			t.logger.Errorf("Error subscribing to topic %s: %s", handler.settings.Topic, token.Error())
			return token.Error()
		}

		t.logger.Infof("Subscribed to topic: %s,broker:%s", handler.settings.Topic, t.settings.Broker)
	}
	url := "tcp://0.0.0.0:" + t.settings.Port
	if err = t.runServer(url); err != nil {
		t.logger.Errorf("runServer url:%s fail,err:%+v", url, err)
	} else {
		t.logger.Infof("runServer listen url:%s", url)
	}
	return err
}

// Stop implements ext.Trigger.Stop
func (t *Trigger) Stop() error {

	//unsubscribe from topics
	for _, handler := range t.handlers {
		t.logger.Debug("Unsubscribing from topic: ", handler.settings.Topic)
		if token := t.client.Unsubscribe(handler.settings.Topic); token.Wait() && token.Error() != nil {
			t.logger.Errorf("Error unsubscribing from topic %s: %s", handler.settings.Topic, token.Error())
		}
	}

	t.client.Disconnect(250)

	return nil
}

func (t *Trigger) getHanlder(handler *clientHandler) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//qos := msg.Qos()
		payload := msg.Payload()

		t.logger.Debugf("Topic[%s] - Payload Recieved: %s", topic, string(payload))

		_, err := runHandler(handler.handler, payload)
		if err != nil {
			t.logger.Error("Error handling message: %v", err)
			return
		}

		//if handler.settings.TopicDown != "" {
		//	reply := &Reply{}
		//	err = reply.FromMap(result)
		//	if err != nil {
		//		t.logger.Error("Error handling message: %v", err)
		//		return
		//	}
		//
		//	if reply.Data != nil {
		//		dataJson, err := json.Marshal(reply.Data)
		//		if err != nil {
		//			return
		//		}
		//		token := client.Publish(handler.settings.ReplyTopic, qos, false, string(dataJson))
		//		sent := token.WaitTimeout(5000 * time.Millisecond)
		//		if !sent {
		//			t.logger.Errorf("Timeout occurred while trying to publish reply to topic '%s'", handler.settings.ReplyTopic)
		//			return
		//		}
		//	}
		//}
	}
}

// RunHandler runs the handler and associated action
func runHandler(handler trigger.Handler, payload []byte) (map[string]interface{}, error) {

	out := &Output{}
	out.Message = payload

	results, err := handler.Handle(context.Background(), out)
	if err != nil {
		return nil, err
	}

	return results, nil
}

//run server
func (t *Trigger) runServer(url string) error {
	server, err := transport.Launch(url)
	if err != nil {
		return err
	}
	backend := broker.NewMemoryBackend()
	backend.SessionQueueSize = 100
	backend.Logger = func(e broker.LogEvent, client *broker.Client, pkt packet.Generic, msg *packet.Message, err error) {
		switch e {
		case broker.NewConnection:
			t.logger.Infof("[%s] new connect event:%s", client.ID(), e)
		case broker.LostConnection:
			t.logger.Infof("[%s] lost connect event:%s", client.ID(), e)
		case broker.LoginConnectSuccess:
			t.logger.Infof("[%s] login success event:%s,", client.ID(), e)
			if msg != nil && msg.Payload != nil {
				if id, thingId, err := parseToken(t.settings.TokenSsl, parseUser(msg.Payload)); err == nil {
					dev := NewDevice(id, thingId, client, t)
					t.deviceCon.Set(client.ID(), dev)
					if err := dev.Notify(DEVICE_STATUS_ONLINE); err != nil {
						t.logger.Errorf("notify device fail id:%s,thingId:%s", id, thingId)
					}
				} else {
					t.logger.Errorf("token check error:%s", err.Error())
					client.Close()
				}
			}
		case broker.ClientDisconnected:
			t.logger.Infof("[%s] client lost event:%s", client.ID(), e)
			if dev := t.deviceCon.Get(client.ID()); dev != nil {
				if err := dev.Notify(DEVICE_STATUS_OFFLINE); err != nil {
					t.logger.Errorf("notify device fail dev:%+v", dev)
				}
			}
		case broker.PacketReceived:
		case broker.MessagePublished:
			if msg != nil && msg.Payload != nil {
				t.logger.Infof("[%s] up send topic:%s,recv:%s", client.ID(), msg.Topic, string(msg.Payload))
				if dev := t.deviceCon.Get(client.ID()); dev != nil {
					if err := dev.Up(msg); err != nil {
						t.logger.Errorf("notify device fail dev:%+v", dev)
					}
				}
			}
		case broker.MessageForwarded:
			if msg != nil && msg.Payload != nil {
				t.logger.Infof("[%s] forwarded topic:%s,recv:%s", client.ID(), msg.Topic, string(msg.Payload))
			}
		default:
			t.logger.Infof("[%s] event:%s", client.ID(), e)
		}
	}
	engine := broker.NewEngine(backend)
	engine.Accept(server)
	return nil
}
