package encodex

import (
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
	Color            string
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
	//event
	event := make(map[string]interface{}, 0)
	event["id"] = uuid.NewV4().String()
	event["version"] = "v1.0.1"
	obj, _ := coerce.ToObject(input.ToMap()["message"])
	objp, _ := coerce.ToObject(obj["params"])
	params, err := buildMessage(objp, a.EventMappings)
	if err != nil {
		ctx.Logger().Errorf("buildMessage fail:%s", err.Error())
		return false, err
	}
	ctx.Logger().Infof("[event] convert params:%+v", params)
	for k, v := range a.AddMappings {
		obj, _ := coerce.ToObject(v)
		obj["time"] = time.Now().Unix() * 1000
		params[k] = v
	}
	event["params"] = params
	if len(params) != 0 && input.Color != a.Color {
		msg := make(map[string]interface{})
		msg["message"] = event
		msg["topic"] = buildUpTopic(a.devices[0].DeviceId, a.devices[0].ThingId, a.EventId)

		ctx.Logger().Infof("[event] topic:%s,encode:%+v", msg["topic"], event)
		output.Data = append(output.Data, msg)
		a.Color = input.Color
	}

	//property
	property := make(map[string]interface{}, 0)
	property["id"] = uuid.NewV4().String()
	property["version"] = "v1.0.1"
	params, err = buildMessage(objp, a.PropertyMappings)
	if err != nil {
		ctx.Logger().Errorf("buildMessage fail:%s", err.Error())
		return false, err
	}
	ctx.Logger().Infof("[property] convert params:%+v", params)
	if len(params) != 0 {
		property["params"] = params
		msgp := make(map[string]interface{})
		msgp["message"] = property
		msgp["topic"] = buildUpPropertyTopic(a.devices[0].DeviceId, a.devices[0].ThingId)

		ctx.Logger().Infof("[property] topic:%s,encode:%+v", msgp["topic"], property)
		output.Data = append(output.Data, msgp)
	}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
