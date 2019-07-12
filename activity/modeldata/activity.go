package modeldata

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/coerce"
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
	ctx.Logger().Infof("modeldata eval:%+v", input)
	message := &ThingMsg{
		Id:      uuid.NewV4().String(),
		Version: "v1.1.0",
		Params:  make(map[string]*ThingData),
	}
	//random
	for _, v := range input.Device {
		if m, err := coerce.ToObject(v); err != nil {
			return false, errors.New("input to obj error")
		} else {
			tmp := &ThingData{
				Id:    m["Id"].(string),
				Type:  m["Type"].(string),
				Value: m["Value"],
				Time:  time.Now().Unix() * 1000,
			}
			message.Params[m["Name"].(string)] = tmp
			ctx.Logger().Infof("eval:%+v", input)
			output := &Output{}
			result, err := json.Marshal(message)
			if err != nil {
				ctx.Logger().Error(err)
				return false, err
			}
			output.Message = string(result)
			output.Topic = buildUpTopic(input.ThingId, input.DeviceId)
			err = ctx.SetOutputObject(output)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}
