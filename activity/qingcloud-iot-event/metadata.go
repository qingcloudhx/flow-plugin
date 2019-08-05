package qingcloud_iot_event

import "github.com/qingcloudhx/core/data/coerce"

type Settings struct {
	Broker       string                 `md:"broker"`       // The broker URL
	Id           string                 `md:"id"`           // The id of client
	Username     string                 `md:"username"`     // The user's name
	Password     string                 `md:"password"`     // The user's password
	Store        string                 `md:"store"`        // The store for message persistence
	CleanSession bool                   `md:"cleanSession"` // Clean session flag
	Qos          int                    `md:"qos"`          // The Quality of Service
	SSLConfig    map[string]interface{} `md:"sslConfig"`    // SSL Configuration
	EventId      string                 `md:"eventId"`
}

type Input struct {
	Device map[string]interface{} `md:"device"` // The message to send
}

type Output struct {
	Data interface{} `md:"data"` // The data recieved
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"device": i.Device,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	i.Device, _ = coerce.ToObject(values["device"])
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Data = values["data"]
	return nil
}
