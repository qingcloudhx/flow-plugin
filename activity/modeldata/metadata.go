package model

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
import (
	"github.com/pkg/errors"
	"github.com/qingcloudhx/core/data/coerce"
)

type DataModel struct {
	Name  string      `md:"name"`
	Id    string      `md:"id"`
	Type  string      `md:"type"`
	Value interface{} `md:"value"`
}
type Input struct {
	DeviceId string      `md:"deviceId"`
	ThingId  string      `md:"thingId"`
	Device   []DataModel `md:"device"`
}
type Settings struct {
	DeviceId string `md:"deviceId"`
	ThingId  string `md:"thingId"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"device":   i.Device,
		"deviceId": i.DeviceId,
		"thingId":  i.ThingId,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	if v, ok := values["device"]; ok {
		if val, ok := v.([]DataModel); ok {
			i.Device = val
		} else {
			err = errors.New("input formMap error")
		}
	} else {
		err = errors.New("input formMap device error")
	}
	if v, ok := values["deviceId"]; ok {
		i.DeviceId, err = coerce.ToString(v)
	} else {
		err = errors.New("input formMap deviceId error")
	}

	if v, ok := values["thingId"]; ok {
		i.ThingId, err = coerce.ToString(v)
	} else {
		err = errors.New("input formMap ThingId error")
	}
	return err
}

type Output struct {
	Message string `md:"message"`
	Topic   string `md:"topic"`
}

func (i *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.Message,
		"topic":   i.Topic,
	}
}

func (i *Output) FromMap(values map[string]interface{}) error {

	var err error
	i.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}
	i.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}
	return nil
}
