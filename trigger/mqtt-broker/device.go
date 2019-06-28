package mqtt

import (
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
	return device
}
func (d Device) Notify(status string) error {
	return nil
}
func (d *Device) Up(data interface{}) error {
	logger := d.trigger.logger
	switch v := data.(type) {
	case *packet.Message:
		if data, err := decode(v.Payload); err != nil {
			//up message
			if token := d.trigger.client.Publish(v.Topic, byte(v.QOS), v.Retain, data); !token.WaitTimeout(5000 * time.Millisecond) {
				logger.Errorf("topic:%s.,data:%s, Publish error:%+v", v.Topic, string(v.Payload), token.Error())
			}
		}
	default:
		logger.Errorf("device up data type error:%+v", v)
	}
	return nil
}
func (d *Device) Down(data interface{}) error {
	logger := d.trigger.logger
	switch v := data.(type) {
	case *packet.Message:
		if data, err := encode(v.Payload); err != nil {
			v.Payload = data.([]byte)
			publish := packet.NewPublish()
			publish.Message = *v
			if err := d.client.Conn().Send(publish, true); err != nil {
				logger.Errorf("device down data type error:%s", err.Error())
			}
		}
	default:
		logger.Errorf("device down data type error:%+v", v)
	}
	return nil
}
