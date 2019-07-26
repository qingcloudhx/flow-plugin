package qingcloud_mqtts_trigger

import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Settings struct {
	Broker        string                 `md:"broker,required"` // The broker URL
	Id            string                 `md:"id,required"`     // The id of client
	Username      string                 `md:"username"`        // The user's name
	Password      string                 `md:"password"`        // The user's password
	Store         string                 `md:"store"`           // The store for message persistence
	CleanSession  bool                   `md:"cleanSession"`    // Clean session flag
	KeepAlive     int                    `md:"keepAlive"`       // Keep Alive time in seconds
	AutoReconnect bool                   `md:"autoReconnect"`   // Enable Auto-Reconnect
	SSLConfig     map[string]interface{} `md:"sslConfig"`       // SSL Configuration
}

type HandlerSettings struct {
	Topics     []interface{} `md:"topics,required"` // The topic to listen on
	ReplyTopic string        `md:"replyTopic"`      // The topic to reply on
	Qos        int           `md:"qos"`             // The Quality of Service
}

type Output struct {
	Topic   string `md:"topic"`
	Message []byte `md:"message"` // The message recieved
}

type Reply struct {
	Topic   string `md:"topic"`
	Message []byte `md:"message"` // The message recieved
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
		"topic":   o.Topic,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToBytes(values["message"])
	if err != nil {
		return err
	}
	o.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}
	return nil
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": r.Message,
		"topic":   r.Topic,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {
	var err error
	r.Message, err = coerce.ToBytes(values["message"])
	if err != nil {
		return err
	}
	r.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}
	return nil
}
