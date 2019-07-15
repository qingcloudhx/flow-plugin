package qingcloud_log

import (
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/coerce"
)

func init() {
	_ = activity.Register(&Activity{})
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

var activityMd = activity.ToMetadata(&Input{})

// Activity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	ctx.GetInputObject(input)

	ctx.Logger().Infof("head:%+v,body:%+v", input.Head, input.Body)

	return true, nil
}
