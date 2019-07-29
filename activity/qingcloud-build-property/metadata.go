package qingcloud_build_property

import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Settings struct {
	Device map[string]interface{} `md:"device,required"` // The broker URL
}

type Output struct {
	Device map[string]interface{} `md:"device"` // The data pulled from the timer
}
type Input struct {
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"device": o.Device,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Device, err = coerce.ToObject(values["device"])
	return err
}
