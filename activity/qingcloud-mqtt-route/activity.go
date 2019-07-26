package qingcloud_mqtt_route

import (
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {
	ctx.Logger().Infof("activity init")
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		ctx.Logger().Error(err)
		return nil, err
	}

	ctx.Logger().Infof("activity init setting:%+v", settings)
	act := &Activity{settings: settings, log: ctx.Logger()} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	settings *Settings
	log      log.Logger
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
		return true, err
	}
	a.log.Debugf("eval start:%+v", input)
	output := &Output{}
	if v, ok := a.settings.Route[input.Topic]; ok {
		for _, val := range v {
			temp := make(map[string]interface{})
			temp["topic"] = val
			temp["message"] = input.Message
			output.Data = append(output.Data, temp)
		}
	}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	a.log.Debugf("eval end:%+v", output)
	return true, nil
}
