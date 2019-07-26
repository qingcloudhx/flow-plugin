package qingcloud_mqtt_route

import "github.com/qingcloudhx/core/data/coerce"

type Settings struct {
	Route map[string]interface{} `md:"route"`
}

type Input struct {
	Topic   string `md:"topic"`
	Message []byte `md:"message"` // The message recieved
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Message, err = coerce.ToBytes(values["message"])
	if err != nil {
		return err
	}
	i.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.Message,
		"topic":   i.Topic,
	}
}

type Output struct {
	Data []interface{} `md:"data"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Data, err = coerce.ToArray(values["data"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}
