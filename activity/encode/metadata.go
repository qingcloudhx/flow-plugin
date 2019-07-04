package encode

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
import (
	"github.com/qingcloudhx/core/data/coerce"
)

type Input struct {
	Label      string  `md:"label"`      // The error message
	Confidence float64 `md:"confidence"` // The error data
	Image      string  `md:"image"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"label":      i.Label,
		"confidence": i.Confidence,
		"image":      i.Image,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Label, err = coerce.ToString(values["label"])
	if err != nil {
		return err
	}
	i.Confidence, err = coerce.ToFloat64(values["confidence"])
	if err != nil {
		return err
	}
	i.Image, err = coerce.ToString(values["image"])
	if err != nil {
		return err
	}
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
