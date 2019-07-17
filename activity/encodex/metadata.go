package encodex

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Input struct {
	Color   string
	Message map[string]interface{}
}
type Settings struct {
	Devices          []interface{}          `md:"devices"`
	EventId          string                 `md:"eventId"`
	Version          string                 `md:"version"`
	EventType        string                 `md:"type"`
	EventMappings    map[string]interface{} `md:"eventMappings"`
	PropertyMappings map[string]interface{} `md:"propertyMappings"`
	AddMappings      map[string]interface{} `md:"addMappings"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.Message,
		"color":   i.Color,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	if v, err := coerce.ToObject(values["message"]); err != nil {
		return err
	} else {
		i.Message = v
	}
	if v, err := coerce.ToString(values["color"]); err != nil {
		return err
	} else {
		i.Color = v
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
