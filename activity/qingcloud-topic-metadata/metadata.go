package qingcloud_topic_metadata

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
	Head map[string]interface{} `md:"head,required"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	head, _ := coerce.ToObject(values["head"])
	i.Head = head
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"head": i.Head,
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
