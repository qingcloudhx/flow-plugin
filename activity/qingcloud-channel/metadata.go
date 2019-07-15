package qingcloud_channel

import "github.com/qingcloudhx/core/data/coerce"

type Settings struct {
	Event string `md:"event,required"` // The broker URL
}

type Input struct {
	Head map[string]interface{} `md:"head"`
	Body []byte                 `md:"body"` // The message to log
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"head": i.Head,
		"body": i.Body,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Head, err = coerce.ToObject(values["head"])
	if err != nil {
		return err
	}
	i.Body, err = coerce.ToBytes(values["body"])
	if err != nil {
		return err
	}

	return nil
}
