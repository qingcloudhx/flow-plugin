package mdmp

import "github.com/qingcloudhx/core/data/coerce"

/**
* @Author: hexing
* @Date: 19-6-27 上午10:54
 */
type Settings struct {
	//Address string `md:"address,required"` // The network type
	//GroupId string `md:"GroupId"`          // Data delimiter for read and write
	//Topic   int    `md:"topic,required"`
}

type HandlerSettings struct {
}

type Output struct {
	Data string `md:"data"` // The data received from the connection
}

type Reply struct {
	Reply string `md:"reply"` // The reply to be sent back
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Data, err = coerce.ToString(values["data"])
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
