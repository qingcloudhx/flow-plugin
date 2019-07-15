package mqttbroker

import (
	"context"
	"sync"
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
	Up(data *Output) error
	Down(data *Output) error
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
}

//build
func NewDevice(id string, trigger *Trigger) DeviceHandler {
	device := &Device{
		id:      id,
		trigger: trigger,
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
func (d *Device) Down(data *Output) error {
	return nil
}
