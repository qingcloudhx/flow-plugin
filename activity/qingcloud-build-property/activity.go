package qingcloud_build_property

import (
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
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
	ctx.Logger().Infof("eval settings:%+v", settings.Device)
	act := &Activity{settings: settings, log: ctx.Logger()}
	return act, nil
}

// Activity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
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
		ctx.Logger().Error(err)
		return false, nil
	}

	output := build(a.settings)
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	a.log.Infof("eval end:%+v", output)
	return true, nil
}
