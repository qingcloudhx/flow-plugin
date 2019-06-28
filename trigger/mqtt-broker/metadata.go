package mqtt

import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Settings struct {
	url string `md:"broker,required"` // The broker URL
}

type HandlerSettings struct {
	TopicDown string `md:"topic_down,required"` // The topic to listen on
	TopicUp   string `md:"topic_up,required"`   // The topic to reply on
	Qos       int    `md:"qos"`                 // The Quality of Service
}

type Output struct {
	Message string `md:"message"` // The message recieved
}

type Reply struct {
	Data interface{} `md:"data"` // The data to reply with
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	return nil
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": r.Data,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	r.Data = values["data"]
	return nil
}
