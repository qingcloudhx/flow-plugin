package qingcloud_tcp

import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Settings struct {
	Network   string `md:"network"`       // The network type
	Host      string `md:"host"`          // The host name or IP for TCP server.
	Port      string `md:"port,required"` // The port to listen on
	Delimiter string `md:"delimiter"`     // Data delimiter for read and write
	TimeOut   int    `md:"timeout"`
}

type HandlerSettings struct {
}

type Output struct {
	Head map[string]interface{} `md:"head"`
	Body []byte                 `md:"body"`
}

type Reply struct {
	Reply []byte `md:"reply"` // The reply to be sent back
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
	r.Reply, err = coerce.ToBytes(values["reply"])
	if err != nil {
		return err
	}

	return nil
}
