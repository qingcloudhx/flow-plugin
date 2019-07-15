package qingcloud_mqtt

import (
	"errors"
	"fmt"
	"github.com/256dpi/gomqtt/broker"
	"github.com/256dpi/gomqtt/packet"
	"github.com/256dpi/gomqtt/transport"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/engine/channels"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/trigger"
	"runtime"
	"strings"
)

/**
* @Author: hexing
* @Date: 19-6-27 上午10:54
 */
var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Factory is a kafka trigger factory
type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
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

// Trigger is a mdmp trigger
type Trigger struct {
	settings  *Settings
	handlers  []trigger.Handler
	logger    log.Logger
	deviceCon *DeviceCon
	//Consumer client.Client //consumer data
}

// Initialize initializes the trigger
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.handlers = ctx.GetHandlers()
	t.logger = ctx.Logger()
	//server
	t.deviceCon = NewDeviceCon()
	return nil
}

// Start starts the kafka trigger
func (t *Trigger) Start() error {
	var err error
	url := t.settings.Url
	if err = t.runServer(url); err != nil {
		t.logger.Errorf("runServer url:%s fail,err:%+v", url, err)
		return err
	} else {
		t.logger.Infof("runServer listen url:%s", url)
	}
	if err = t.register(); err != nil {
		t.logger.Errorf("register fail,err:%+v", err)
	}
	return nil
}

// Start starts the kafka trigger
func (t *Trigger) Stop() error {
	return nil
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
			fallthrough
		case broker.ClientDisconnected:
			t.logger.Infof("[%s] client lost event:%s", client.ID(), e)
			if dev := t.deviceCon.Get(client.ID()); dev != nil {
				data := buildPackage(buildHead(mqtt_cmd_disconnect, client.ID(), "", ""), []byte{})
				if err := dev.Up(data); err != nil {
					t.logger.Errorf("dev up cmd:%s error:%s", mqtt_cmd_connect, err)
				}
				t.deviceCon.Del(client.ID())
			}
		case broker.PacketReceived:
			t.logger.Infof("[%s] client recv event:%s", client.ID(), e)
			if pkt != nil {
				if v, ok := pkt.(*packet.Connect); ok {
					t.logger.Infof("[%s] client create data:%s", v.ClientID, v.String())
					dev := NewDevice(v.ClientID, client, t)
					t.deviceCon.Set(v.ClientID, dev)
					data := buildPackage(buildHead(mqtt_cmd_connect, v.ClientID, v.Username, v.Password), []byte{})
					if err := dev.Up(data); err != nil {
						t.logger.Errorf("dev up cmd;%s error:%s", mqtt_cmd_connect, err)
					}
				}
			}
		case broker.MessagePublished:
			if msg != nil && msg.Payload != nil {
				t.logger.Infof("[%s] up send topic:%s,recv:%s", client.ID(), msg.Topic, string(msg.Payload))
				//id := parseTopic(msg.Topic)
				if dev := t.deviceCon.Get(client.ID()); dev != nil {
					data := buildPackage(buildHead(mqtt_cmd_data, client.ID(), "", ""), msg.Payload)
					if err := dev.Up(data); err != nil {
						t.logger.Errorf("dev up cmd:%s error:%s", mqtt_cmd_data, err)
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
func (t *Trigger) onMessage(msg interface{}) {
	defer func() {
		if err := recover(); err != nil {
			PrintStack()
		}
	}()
	if v, ok := msg.(map[string]interface{}); ok {
		if head, ok := v[mqtt_head]; ok {
			if id, ok := head.(map[string]interface{})[mqtt_client_id].(string); ok {
				if dev := t.deviceCon.Get(id); dev != nil {
					err := dev.Down(msg)
					if err != nil {
						t.logger.Errorf("down error:%s", err.Error())
					}
				}
			}
		}
	} else {
		t.logger.Infof("result:%+v", msg)
	}
}

//Register event
func (t *Trigger) register() error {
	event := t.settings.Event
	t.logger.Infof("Started register event %s", event)
	if event != "" {
		e := strings.Split(event, ",")
		for _, v := range e {
			ch := channels.Get(v)
			if ch == nil {
				return errors.New(fmt.Sprintf("channels:%s not existed", v))
			}
			err := ch.RegisterCallback(t.onMessage)
			if err != nil {
				return errors.New(fmt.Sprintf("RegisterCallback error:%s", err.Error()))
			}
		}
	}
	return nil
}
func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	log.RootLogger().Errorf("panic ==> %s\n", string(buf[:n]))
}
