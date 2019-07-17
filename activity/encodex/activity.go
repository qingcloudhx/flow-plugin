package encodex

import (
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
	act := &Activity{devices: devices, EventId: settings.EventId, Version: settings.Version, EventMappings: settings.EventMappings, PropertyMappings: settings.PropertyMappings, AddMappings: settings.AddMappings}
	return act, nil
}

// Activity is an Activity that used to cause an explicit error in the flow
// inputs : {message,data}
// outputs: node
type Activity struct {
	devices          []*DeviceInfo
	EventId          string `json:"eventId"`
	Version          string `json:"version"`
	EventMappings    map[string]interface{}
	AddMappings      map[string]interface{}
	PropertyMappings map[string]interface{}
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
	message := make(map[string]interface{}, 0)
	message["id"] = uuid.NewV4().String()
	message["version"] = "v1.0.1"
	params, err := buildMessage(input.ToMap(), a.EventMappings)
	if err != nil {
		ctx.Logger().Errorf("buildMessage fail:%s", err.Error())
		return false, err
	}
	message["params"] = params
	for k, v := range a.AddMappings {
		params[k] = v
	}
	//data, _ := json.Marshal(message)
	msg := make(map[string]interface{})
	msg["message"] = message
	msg["topic"] = buildUpTopic(a.devices[0].DeviceId, a.devices[0].ThingId, a.EventId)
	ctx.Logger().Infof("[event] topic:%s,encode:%+v", msg["topic"], message)
	output.Data = append(output.Data, msg)

	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
