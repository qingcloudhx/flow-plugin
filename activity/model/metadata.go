package model

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Input struct {
}
type Settings struct {
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	return nil
}

type Output struct {
	Message string `md:"message"`
}

func (i *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.Message,
	}
}

func (i *Output) FromMap(values map[string]interface{}) error {

	var err error
	i.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}
	return nil
}
