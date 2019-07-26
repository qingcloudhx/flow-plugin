package qingcloud_mqtts_trigger

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/support/ssl"
	"github.com/qingcloudhx/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})
var flag int32 = 0

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Trigger is simple MQTT trigger
type Trigger struct {
	handlers map[string]*clientHandler
	settings *Settings
	logger   log.Logger
	options  *mqtt.ClientOptions
	client   mqtt.Client
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

		for _, v := range s.Topics {
			t.handlers[v.(string)] = &clientHandler{handler: handler, settings: s}
		}
	}

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
func (t *Trigger) Start() error {

	client := mqtt.NewClient(t.options)

	if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		t.logger.Error(token.Error())
		go func() {
			for {
				if token := client.Connect(); token.Wait() && token.Error() != nil {
					t.logger.Error(token.Error())
					time.Sleep(3 * time.Second)
				} else {
					t.logger.Infof("mqtt conect success")
					atomic.StoreInt32(&flag, 1)
					break
				}
			}
		}()
	} else {
		t.logger.Infof("mqtt conect success")
		atomic.StoreInt32(&flag, 1)
	}
	t.client = client

	for k, handler := range t.handlers {
		for _, v := range handler.settings.Topics {
			val, _ := v.(string)
			if k == val {
				if token := client.Subscribe(val, byte(handler.settings.Qos), t.getHanlder(handler)); token.Wait() && token.Error() != nil {
					t.logger.Errorf("Error subscribing to topic %s: %s", val, token.Error())
					return token.Error()
				}
			}
		}
		t.logger.Debugf("Subscribed to topic: %+v", handler.settings.Topics)
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *Trigger) Stop() error {

	//unsubscribe from topics
	for k, handler := range t.handlers {
		t.logger.Debug("Unsubscribing from topic: ", handler.settings.Topics)
		for _, v := range handler.settings.Topics {
			val, _ := v.(string)
			if k == val {
				if token := t.client.Unsubscribe(val); token.Wait() && token.Error() != nil {
					t.logger.Errorf("Error unsubscribing from topic %s: %s", v, token.Error())
				}
			}
		}
	}

	t.client.Disconnect(250)

	return nil
}

func (t *Trigger) getHanlder(handler *clientHandler) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		qos := msg.Qos()
		payload := msg.Payload()

		t.logger.Debugf("Topic[%s] - Payload Recieved: %s", topic, payload)

		result, err := runHandler(handler.handler, topic, payload)
		if err != nil {
			t.logger.Error("Error handling message: %v", err)
			return
		}

		if handler.settings.ReplyTopic != "" {
			reply := &Reply{}
			err = reply.FromMap(result)
			if err != nil {
				t.logger.Error("Error handling message: %v", err)
				return
			}

			if reply.Message != nil && reply.Topic != "" {
				token := client.Publish(reply.Topic, qos, false, reply.Message)
				sent := token.WaitTimeout(5000 * time.Millisecond)
				if !sent {
					t.logger.Errorf("Timeout occurred while trying to publish reply to topic '%s'", handler.settings.ReplyTopic)
					return
				}
			}
		}
	}
}

// RunHandler runs the handler and associated action
func runHandler(handler trigger.Handler, topic string, payload []byte) (map[string]interface{}, error) {

	out := &Output{}
	out.Topic = topic
	out.Message = payload

	results, err := handler.Handle(context.Background(), out)
	if err != nil {
		return nil, err
	}

	return results, nil
}
