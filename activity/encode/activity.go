package encode

import (
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
	act := &Activity{devices: devices, EventId: settings.EventId, Version: settings.Version, EventMappings: settings.EventMappings, PrppertyMappings: settings.PropertyMappings}
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
	PrppertyMappings map[string]interface{}
	Type             string
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
	output.Type = input.Type
	output.Data = make([]interface{}, 0)
	switch input.Type {
	case "heartbeat":
		message := buildHeartBeat(a.devices[id].DeviceId, a.devices[id].ThingId, DEVICE_STATUS_ONLINE)
		msg := make(map[string]interface{})
		msg["message"] = string(message)
		msg["topic"] = buildHeartbeatTopic(a.devices[id].DeviceId, a.devices[id].ThingId)
		output.Data = append(output.Data, msg)
	case "data":
		if exists(input.License) {
			ctx.Logger().Infof("filter repeat license:%s", input.License)
			msg := make(map[string]interface{})
			msg["message"] = ""
			msg["topic"] = ""
			output.Data = append(output.Data, msg)
		} else {
			add(input.License, 3*time.Second, input.License)

			message := make(map[string]interface{}, 0)
			message["id"] = uuid.NewV4().String()
			message["version"] = "v1.0.1"
			params, err := buildMessage(input.ToMap(), a.EventMappings)
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

			//data, _ := json.Marshal(message)
			msg := make(map[string]interface{})
			msg["message"] = message
			msg["topic"] = buildUpTopic(a.devices[id].DeviceId, a.devices[id].ThingId, a.EventId)
			ctx.Logger().Infof("[event] topic:%s,encode:%+v", msg["topic"], message)
			output.Data = append(output.Data, msg)

			//todo property
			message["id"] = uuid.NewV4().String()
			message["version"] = "v1.0.1"
			params = make(map[string]interface{})
			params, err = buildMessage(input.ToMap(), a.PrppertyMappings)
			if err != nil {
				ctx.Logger().Errorf("buildMessage fail:%s", err.Error())
				return false, err
			}
			msg["message"] = message
			msg["topic"] = buildUpPropertyTopic(a.devices[id].DeviceId, a.devices[id].ThingId, a.EventId)
			ctx.Logger().Infof("[property] topic:%s,encode:%+v", msg["topic"], message)
			output.Data = append(output.Data, msg)
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
