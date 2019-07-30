package delay

import "github.com/qingcloudhx/core/data/coerce"

/**
* @Author: hexing
* @Date: 19-7-30 下午1:26
 */
type Input struct {
	Data []interface{} `md:"data"`
	Type string        `md:"type"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": i.Data,
		"type": i.Type,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Data, err = coerce.ToArray(values["data"])
	if err != nil {
		return err
	}
	i.Type, err = coerce.ToString(values["type"])
	if err != nil {
		return err
	}
	return nil
}

type Output struct {
	Data []interface{} `md:"data"`
	Type string        `md:"type"`
}

func (i *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": i.Data,
		"type": i.Type,
	}
}

func (i *Output) FromMap(values map[string]interface{}) error {
	var err error
	i.Data, err = coerce.ToArray(values["data"])
	if err != nil {
		return err
	}
	i.Type, err = coerce.ToString(values["type"])
	if err != nil {
		return err
	}
	return nil
}
