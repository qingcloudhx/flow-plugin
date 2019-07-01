package mqtt_broker

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"time"

	"github.com/qingcloudhx/gomqtt/broker"
	"github.com/qingcloudhx/gomqtt/packet"
)

/**
* @Author: hexing
* @Date: 19-6-28 下午3:23
 */
const (
	IOT_DEVICE_STATUS_END  = "iote-global-onoffline-end"
	IOT_DEVICE_STATUS_EDGE = "iote-global-onoffline-edge"
	DEVICE_STATUS_ONLINE   = "online"  // 在线
	DEVICE_STATUS_OFFLINE  = "offline" // 离线
)

type DeviceHandler interface {
	Notify(status string) error
	Up(data interface{}) error
	Down(data interface{}) error
}
type DeviceCon struct {
	Channels map[string]DeviceHandler
	lock     sync.RWMutex
}

//build map
func NewDeviceCon() *DeviceCon {
	c := &DeviceCon{
		Channels: make(map[string]DeviceHandler),
	}
	return c
}
func (dev *DeviceCon) Set(id string, handler DeviceHandler) {
	dev.lock.RLock()
	defer dev.lock.RUnlock()
	dev.Channels[id] = handler
	return
}
func (dev *DeviceCon) Get(id string) DeviceHandler {
	dev.lock.RLock()
	defer dev.lock.RUnlock()
	if v, ok := dev.Channels[id]; ok {
		return v
	}
	return nil
}
func (dev *DeviceCon) Del(id string) DeviceHandler {
	dev.lock.RLock()
	defer dev.lock.RUnlock()
	delete(dev.Channels, id)
	return nil
}

type Device struct {
	id      string
	thingId string
	client  *broker.Client
	trigger *Trigger
}

//build
func NewDevice(id, thingId string, client *broker.Client, trigger *Trigger) DeviceHandler {
	device := &Device{
		id:      id,
		thingId: thingId,
		client:  client,
		trigger: trigger,
	}
	//reg up topic

	return device
}

//hearbeat notify
func (d *Device) Notify(status string) error {
	logger := d.trigger.logger
	logger.Infof("[status] notify id:%s,thingId:%s,status:%s", d.id, d.thingId, status)
	data := buildHeartBeat(d.id, d.thingId, status)
	//up message
	topic := buildStatusToipc(d.id, d.thingId)
	if token := d.trigger.client.Publish(topic, byte(0), false, data); !token.WaitTimeout(5000 * time.Millisecond) {
		logger.Errorf("[status] topic:%s.,data:%s, Publish error:%+v", topic, string(data), token.Error())
	} else {
		logger.Infof("[status] notify up id:%s,topic:%s,data:%s", d.id, topic, string(data))
	}
	if status == DEVICE_STATUS_ONLINE {
		topic = buildDownTopic(d.id, d.thingId)
		if token := d.trigger.client.Subscribe(topic, byte(0), callHanlder(d)); token.Wait() && token.Error() != nil {
			logger.Errorf("Error subscribing to topic %s: %s", topic, token.Error())
			return token.Error()
		}
	} else {
		topic = buildDownTopic(d.id, d.thingId)
		if token := d.trigger.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
			logger.Errorf("Error unsubscribing to topic %s: %s", topic, token.Error())
			return token.Error()
		}
	}
	return nil
}

//up msg
func (d *Device) Up(data interface{}) error {
	logger := d.trigger.logger
	switch v := data.(type) {
	case *packet.Message:
		if data, err := decode(v.Payload); err == nil {
			//up message
			if token := d.trigger.client.Publish(v.Topic, byte(v.QOS), v.Retain, data); !token.WaitTimeout(5000 * time.Millisecond) {
				logger.Errorf("[up] topic:%s.,data:%s, Publish error:%+v", v.Topic, string(v.Payload), token.Error())
			} else {
				logger.Infof("[up] id:%s,topic:%s", d.id, v.Topic)
			}
		}
	default:
		logger.Errorf("device up data type error:%+v", v)
	}
	return nil
}

//down msg
func (d *Device) Down(data interface{}) error {
	logger := d.trigger.logger
	logger.Infof("[down] id:%s,thing:%s", d.id, d.thingId)
	switch v := data.(type) {
	case mqtt.Message:
		if data, err := encode(v.Payload()); err == nil {
			msg := packet.Message{
				Payload: data.([]byte),
				Topic:   v.Topic(),
				QOS:     packet.QOS(v.Qos()),
				Retain:  v.Retained(),
			}
			pubMsg := packet.NewPublish()
			pubMsg.Message = msg
			pubMsg.ID = packet.ID(uuid())
			if err := d.client.Conn().Send(pubMsg, false); err != nil {
				logger.Errorf("[down] client send fail:%s", err.Error())
			}
		} else {
			logger.Errorf("[down] data encode error:%s", err.Error())
		}
	default:
		logger.Errorf("device down data type error:%+v", v)
	}
	return nil
}

//message hook
func callHanlder(dev *Device) func(mqtt.Client, mqtt.Message) {
	logger := dev.trigger.logger
	return func(client mqtt.Client, msg mqtt.Message) {
		logger.Infof("[down] topic:%s,data:%s", msg.Topic(), string(msg.Payload()))
		if err := dev.Down(msg); err != nil {
			logger.Errorf("[down] id:%s,thingId:%s down fail:%s", err.Error())
		}
	}
}
