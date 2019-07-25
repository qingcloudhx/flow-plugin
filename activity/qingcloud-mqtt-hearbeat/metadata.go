package qingcloud_mqtt_hearbeat

import "github.com/qingcloudhx/core/data/coerce"

type Settings struct {
	Format       string                 `md:"format,required"`
	Broker       string                 `md:"broker,required"` // The broker URL
	Id           string                 `md:"id,required"`     // The id of client
	Username     string                 `md:"username"`        // The user's name
	Password     string                 `md:"password"`        // The user's password
	Store        string                 `md:"store"`           // The store for message persistence
	CleanSession bool                   `md:"cleanSession"`    // Clean session flag
	Qos          int                    `md:"qos"`             // The Quality of Service
	SSLConfig    map[string]interface{} `md:"sslConfig"`       // SSL Configuration
}

type Input struct {
	ThingId  string `md:"thingId,required"`
	DeviceId string `md:"deviceId,required"`
	Status   string `md:"status,required"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	thingId, _ := coerce.ToString(values["thingId"])
	i.ThingId = thingId
	deviceId, _ := coerce.ToString(values["deviceId"])
	i.DeviceId = deviceId
	status, _ := coerce.ToString(values["status"])
	i.DeviceId = status
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"thingId":  i.ThingId,
		"deviceId": i.DeviceId,
		"status":   i.Status,
	}
}

type Output struct {
	Topic   string `md:"topic"`
	Message string `md:"message"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	topic, _ := coerce.ToString(values["topic"])
	o.Topic = topic
	message, _ := coerce.ToString(values["message"])
	o.Message = message
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"topic":   o.Topic,
		"message": o.Message,
	}
}
