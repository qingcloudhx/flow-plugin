package encode

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Input struct {
	Type       string  `md:"type"`
	Id         int     `md:"id"`
	Label      string  `md:"label"`      // The error message
	Confidence float64 `md:"confidence"` // The error data
	Image      string  `md:"image"`
}
type Settings struct {
	Devices []interface{} `md:"devices"`
	EventId string        `md:"eventId"`
	Version string        `md:"version"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"label":      i.Label,
		"confidence": i.Confidence,
		"image":      i.Image,
		"type":       i.Type,
		"id":         i.Id,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	if v, ok := values["label"]; ok {
		i.Label, err = coerce.ToString(v)
		if err != nil {
			return err
		}
	} else {
		i.Label = ""
	}
	if v, ok := values["confidence"]; ok {
		i.Confidence, err = coerce.ToFloat64(v)
		if err != nil {
			return err
		}
	} else {
		i.Confidence = 0
	}
	if v, ok := values["image"]; ok {
		i.Image, err = coerce.ToString(v)
		if err != nil {
			return err
		}
	} else {
		i.Image = ""
	}
	if v, ok := values["type"]; ok {
		i.Type, err = coerce.ToString(v)
		if err != nil {
			return err
		}
	} else {
		i.Type = ""
	}
	if v, ok := values["id"]; ok {
		i.Id, err = coerce.ToInt(v)
		if err != nil {
			return err
		}
	} else {
		i.Id = 0
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
