package mqttbroker

import (
	"context"
	"github.com/256dpi/gomqtt/broker"
	"github.com/256dpi/gomqtt/packet"
	"github.com/qingcloudhx/core/data/coerce"
	"sync"
)

/**
* @Author: hexing
* @Date: 19-6-28 下午3:23
 */
type DeviceHandler interface {
	Up(data *Output) error
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
	trigger *Trigger
	client  *broker.Client
	inc     packet.ID
}

//build
func NewDevice(id string, client *broker.Client, trigger *Trigger) DeviceHandler {
	device := &Device{
		id:      id,
		trigger: trigger,
		client:  client,
	}
	//reg up topic

	return device
}

//up msg
func (d *Device) Up(data *Output) error {
	logger := d.trigger.logger
	for _, handler := range d.trigger.handlers {
		if result, err := handler.Handle(context.Background(), data); err != nil {
			logger.Errorf("handler error:%s", err.Error())
		} else {
			logger.Infof("handler result:%+v", result)
		}
	}
	return nil
}

//down msg
func (d *Device) Down(data interface{}) error {
	message, err := coerce.ToObject(data)
	if err != nil {
		return err
	}
	head, err := coerce.ToObject(message[mqtt_head])
	if err != nil {
		return err
	}
	topic, err := coerce.ToString(head[mqtt_topic])
	if err != nil {
		return err
	}
	Qos, err := coerce.ToInt(head[mqtt_qos])
	if err != nil {
		return err
	}
	payload, err := coerce.ToBytes(message[mqtt_body])
	if err != nil {
		return err
	}
	reply := packet.NewPublish()
	reply.Message.Topic = topic
	reply.Message.Payload = payload
	reply.Message.QOS = packet.QOS(Qos)
	reply.ID = d.inc
	reply.Dup = true
	d.inc++
	return d.client.Conn().Send(reply, false)
}
