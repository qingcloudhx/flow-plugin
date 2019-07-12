package modeldata

import (
	"encoding/json"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	uuid "github.com/satori/go.uuid"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
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
	act := &Activity{}
	return act, nil
}

// Activity is an Activity that used to cause an explicit error in the flow
// inputs : {message,data}
// outputs: node
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval returns an error
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}
	message := &ThingMsg{
		Id:      uuid.NewV4().String(),
		Version: "v1.1.0",
		Params:  make(map[string]*ThingData),
	}
	//random
	for _, v := range input.Device {
		tmp := &ThingData{
			Id:    v.Id,
			Type:  v.Type,
			Value: v.Value,
			Time:  time.Now().Unix() * 1000,
		}
		message.Params[v.Name] = tmp
	}
	ctx.Logger().Infof("eval:%+v", input)
	output := &Output{}
	result, err := json.Marshal(message)
	if err != nil {
		ctx.Logger().Error(err)
		return false, err
	}
	output.Message = string(result)
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
