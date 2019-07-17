package encodex

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/coerce"
	"github.com/qingcloudhx/core/data/metadata"
	uuid "github.com/satori/go.uuid"
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
	var devices []*DeviceInfo
	for _, v := range settings.Devices {
		dev := &DeviceInfo{}
		m, _ := coerce.ToObject(v)
		dev.DeviceId = m["deviceId"].(string)
		dev.ThingId = m["thingId"].(string)
		devices = append(devices, dev)
	}
	act := &Activity{devices: devices, EventId: settings.EventId, Version: settings.Version, Mappings: settings.Mappings, Type: settings.EventType}
	return act, nil
}

// Activity is an Activity that used to cause an explicit error in the flow
// inputs : {message,data}
// outputs: node
type Activity struct {
	devices  []*DeviceInfo
	EventId  string `json:"eventId"`
	Version  string `json:"version"`
	Mappings map[string]interface{}
	Type     string
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
	id := input.Id
	if id > len(a.devices)-1 {
		ctx.Logger().Errorf("eval:%d", len(a.devices))
		return false, errors.New("param error")
	}
	message := make(map[string]interface{}, 0)
	message["id"] = uuid.NewV4().String()
	message["version"] = "v1.0.1"
	params, err := buildMessage(input.Message, a.Mappings)
	if err != nil {
		ctx.Logger().Errorf("buildMessage fail:%s", err.Error())
		return false, err
	}
	message["params"] = params
	data, _ := json.Marshal(message)
	output.Message = string(data)
	if a.Type == "property" {
		output.Topic = buildUpPropertyTopic(a.devices[id].DeviceId, a.devices[id].ThingId, a.EventId)
	} else {
		output.Topic = buildUpTopic(a.devices[id].DeviceId, a.devices[id].ThingId, a.EventId)
	}
	ctx.Logger().Infof("topic:%s,encode:%s", output.Topic, data)
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
