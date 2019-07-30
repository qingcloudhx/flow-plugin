package mqtts

import "github.com/qingcloudhx/core/data/coerce"

type Settings struct {
	Broker       string                 `md:"broker,required"` // The broker URL
	Id           string                 `md:"id,required"`     // The id of client
	Username     string                 `md:"username"`        // The user's name
	Password     string                 `md:"password"`        // The user's password
	Store        string                 `md:"store"`           // The store for message persistence
	CleanSession bool                   `md:"cleanSession"`    // Clean session flag
	Topic        string                 `md:"topic"`           // The topic to publish to
	Qos          int                    `md:"qos"`             // The Quality of Service
	SSLConfig    map[string]interface{} `md:"sslConfig"`       // SSL Configuration
	Delay        int                    `md:"delay"`
}
type DeviceUpStatusMsg struct {
	Status     string `json:"status"`
	UserId     string `json:"user_id"`
	ThingId    string `json:"thing_id"`
	PropertyId string `json:"property_id"`
	DeviceId   string `json:"device_id"`
	Time       int64  `json:"time"`
}
type Input struct {
	Data []interface{} `md:"data"`
	Type string        `md:"type"`
}

type Output struct {
	Data interface{} `md:"data"` // The data recieved
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": i.Data,
		"type": i.Type,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	i.Data, _ = coerce.ToArray(values["data"])
	i.Type, _ = coerce.ToString(values["type"])
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
