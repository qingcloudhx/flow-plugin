package qingcloud_channel

import (
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/engine/channels"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Input{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}
	act := &Activity{settings: settings}
	return act, nil
}

// Activity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type Activity struct {
	settings *Settings
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		ctx.Logger().Error(err)
		return false, nil
	}

	if ch := channels.Get(a.settings.Event); ch != nil {
		data := make(map[string]interface{})
		data["body"] = input.Body
		head := make(map[string]interface{})
		head["topic"] = "test"
		for k, v := range input.Head {
			head[k] = v
		}
		data["head"] = head
		if head["cmd"] == "data" {
			ch.Publish(data)
			ctx.Logger().Infof("channel publish head:%+v,body:%+v", input.Head, input.Body)
		}
	} else {
		ctx.Logger().Infof("not find channel head:%+v,body:%+v", input.Head, input.Body)
	}

	return true, nil
}
