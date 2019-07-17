package encode

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
	var devices []*DeviceInfo
	for _, v := range settings.Devices {
		dev := &DeviceInfo{}
		m, _ := coerce.ToObject(v)
		dev.DeviceId = m["deviceId"].(string)
		dev.ThingId = m["thingId"].(string)
		devices = append(devices, dev)
	}
	act := &Activity{devices: devices, EventId: settings.EventId, Version: settings.Version, Mappings: settings.Mappings}
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
	switch input.Type {
	case "heartbeat":
		message := buildHeartBeat(a.devices[id].DeviceId, a.devices[id].ThingId, DEVICE_STATUS_ONLINE)
		output.Message = string(message)
		output.Topic = buildHeartbeatTopic(a.devices[id].DeviceId, a.devices[id].ThingId)
	case "data":
		if exists(input.License) {
			ctx.Logger().Infof("filter repeat license:%s", input.License)
		} else {
			add(input.License, 3*time.Second, input.License)
			//message := &ThingMsg{
			//	Id:      uuid.NewV4().String(),
			//	Version: a.Version,
			//	Params:  make(map[string]*ThingData),
			//}
			message := make(map[string]interface{}, 0)
			message["id"] = uuid.NewV4().String()
			message["version"] = "v1.0.1"
			params, err := buildMessage(input.ToMap(), a.Mappings)
			if err != nil {
				ctx.Logger().Errorf("buildMessage fail:%s", err.Error())
				return false, err
			}
			//todo fix image
			message["params"] = params
			if image, ok := params["image"]; ok {
				if v, ok := image.(map[string]interface{}); ok {
					if url := getPictureUrl(v["value"].(string), ctx.Logger()); url == "" {
						ctx.Logger().Errorf("getPictureUrl image:%s", input.Image)
						return false, nil
					} else {
						v["value"] = url
					}
				}
			}
			//label := &ThingData{
			//	Id:    "35",
			//	Type:  "string",
			//	Value: input.Label,
			//	Time:  time.Now().Unix() * 1000,
			//}
			//message.Params["label"] = label
			//image := &ThingData{
			//	Id:    "34",
			//	Type:  "string",
			//	Value: input.Image,
			//	Time:  time.Now().Unix() * 1000,
			//}
			//if url := getPictureUrl(input.Image, ctx.Logger()); url == "" {
			//	ctx.Logger().Errorf("getPictureUrl image:%s", input.Image)
			//	return false, nil
			//} else {
			//	image.Value = url
			//}
			//message.Params["image"] = image
			//confidence := &ThingData{
			//	Id:    "36",
			//	Type:  "float",
			//	Value: input.Confidence,
			//	Time:  time.Now().Unix() * 1000,
			//}
			//message.Params["confidence"] = confidence
			//
			//color := &ThingData{
			//	Id:    "60",
			//	Type:  "string",
			//	Value: input.Color,
			//	Time:  time.Now().Unix() * 1000,
			//}
			//message.Params["color"] = color
			//license := &ThingData{
			//	Id:    "61",
			//	Type:  "string",
			//	Value: input.License,
			//	Time:  time.Now().Unix() * 1000,
			//}
			//message.Params["license"] = license

			data, _ := json.Marshal(message)
			output.Message = string(data)
			output.Topic = buildUpTopic(a.devices[id].DeviceId, a.devices[id].ThingId, a.EventId)
			ctx.Logger().Infof("encode:%s", data)
		}
	default:
		ctx.Logger().Errorf("data error:%+v", input)
		return false, nil
	}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
