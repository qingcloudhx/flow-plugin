package encodex

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Input struct {
	Id      int
	Message map[string]interface{}
}
type Settings struct {
	Devices   []interface{}          `md:"devices"`
	EventId   string                 `md:"eventId"`
	Version   string                 `md:"version"`
	EventType string                 `md:"type"`
	Mappings  map[string]interface{} `md:"mappings"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.Message,
		"id":      i.Id,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	if v, err := coerce.ToObject(values["message"]); err != nil {
		return err
	} else {
		i.Message = v
	}
	if v, err := coerce.ToInt(values["id"]); err != nil {
		return err
	} else {
		i.Id = v
	}
	return nil
}

type Output struct {
	Message string `md:"message"`
	Topic   string `md:"topic"`
}

func (i *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.Message,
		"topic":   i.Topic,
	}
}

func (i *Output) FromMap(values map[string]interface{}) error {

	var err error
	i.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}
	i.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}
	return nil
}
