package encode

import (
	"encoding/json"
	"github.com/qingcloudhx/core/activity"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-4 上午11:29
 */
func init() {
	_ = activity.Register(&Activity{})
}

var activityMd = activity.ToMetadata(&Input{})

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
	ctx.Logger().Infof("eval:%+v", input)
	output := &Output{}

	message := &ThingEventMsg{
		Id:      "xxx",
		Version: "v1.0.0",
		Params: &ThingEventData{
			Id:   "iote-a64015b1-5c4d-4ff6-89fe-cca8aed35067",
			Time: time.Now().Unix(),
			Data: make([]*EventData, 0, 3),
		},
	}
	event := &EventData{
		Id:    35,
		Name:  "label",
		Type:  "string",
		Value: input.Label,
	}
	message.Params.Data = append(message.Params.Data, event)
	event = &EventData{
		Id:    34,
		Name:  "image",
		Type:  "string",
		Value: input.Image,
	}
	message.Params.Data = append(message.Params.Data, event)
	event = &EventData{
		Id:    36,
		Name:  "confidence",
		Type:  "float",
		Value: input.Confidence,
	}
	message.Params.Data = append(message.Params.Data, event)
	data, _ := json.Marshal(message)
	output.Message = string(data)
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	ctx.Logger().Infof("encode:%s", data)
	return true, nil
}
