package mqttbroker

import "github.com/qingcloudhx/core/data/coerce"

/**
* @Author: hexing
* @Date: 19-6-27 上午10:54
 */
type Settings struct {
	Url   string `md:"url,required"` // The broker URL
	Event string `md:"event"`
	//Address string `md:"address,required"` // The network type
	//GroupId string `md:"GroupId"`          // Data delimiter for read and write
	//Topic   int    `md:"topic,required"`
}

type HandlerSettings struct {
}

type Output struct {
	Head map[string]interface{} `md:"head"`
	Body []byte                 `md:"body"`
}

type Reply struct {
	Reply string `md:"reply"` // The reply to be sent back
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"head": o.Head,
		"body": o.Body,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Body, err = coerce.ToBytes(values["body"])
	if err != nil {
		return err
	}
	o.Head, err = coerce.ToObject(values["head"])
	if err != nil {
		return err
	}
	return nil
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"reply": r.Reply,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	var err error
	r.Reply, err = coerce.ToString(values["reply"])
	if err != nil {
		return err
	}

	return nil
}
