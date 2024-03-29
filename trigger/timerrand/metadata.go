package timerrand

import (
	"errors"
	"github.com/qingcloudhx/core/data/coerce"
)

type ThingData struct {
	Id    string      `md:"id"`
	Type  string      `md:"type"`
	Value interface{} `md:"value"`
	Name  string      `md:"name"`
}

type Output struct {
	Device   []interface{} `md:"device"` // The data pulled from the timer
	DeviceId string        `md:"deviceId"`
	ThingId  string        `md:"thingId"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"device":   o.Device,
		"deviceId": o.DeviceId,
		"thingId":  o.ThingId,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Device, err = coerce.ToArray(values["device"])
	if v, ok := values["deviceId"]; ok {
		o.DeviceId, err = coerce.ToString(v)
	} else {
		err = errors.New("input formMap deviceId error")
	}

	if v, ok := values["thingId"]; ok {
		o.ThingId, err = coerce.ToString(v)
	} else {
		err = errors.New("input formMap ThingId error")
	}
	return err
}
