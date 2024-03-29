package delay

import (
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"time"
)

func init() {
	_ = activity.Register(&Activity{})
}

//type Input struct {
//	Message    string `md:"message"`    // The message to log
//	AddDetails bool   `md:"addDetails"` // Append contextual execution information to the log message
//}
//
//func (i *Input) ToMap() map[string]interface{} {
//	return map[string]interface{}{
//		"message":    i.Message,
//		"addDetails": i.AddDetails,
//	}
//}
//
//func (i *Input) FromMap(values map[string]interface{}) error {
//
//	var err error
//	i.Message, err = coerce.ToString(values["message"])
//	if err != nil {
//		return err
//	}
//	i.AddDetails, err = coerce.ToBool(values["addDetails"])
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

type Settings struct {
	Delay int `md:"delay,required"`
}

var activityMd = activity.ToMetadata(&Input{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}
	act := &Activity{Delay: settings.Delay}
	return act, nil
}

// Activity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type Activity struct {
	Delay int
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
		ctx.Logger().Errorf("input:%+v", input)
		return false, err
	}
	time.Sleep(time.Duration(a.Delay) * time.Millisecond)
	ctx.Logger().Infof("delay:%d data:%+v", a.Delay, input)
	output := &Output{
		Type: input.Type,
		Data: input.Data,
	}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
