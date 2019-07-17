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
	License    string  `md:"license"`
	Color      string  `md:"color"`
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
		"label":      i.Label,
		"confidence": i.Confidence,
		"image":      i.Image,
		"type":       i.Type,
		"id":         i.Id,
		"license":    i.License,
		"color":      i.Color,
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
	if v, ok := values["license"]; ok {
		i.License, err = coerce.ToString(v)
		if err != nil {
			return err
		}
	} else {
		i.License = ""
	}
	if v, ok := values["color"]; ok {
		i.Color, err = coerce.ToString(v)
		if err != nil {
			return err
		}
	} else {
		i.Color = ""
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
