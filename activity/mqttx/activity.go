package mqttx

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/coerce"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/support/ssl"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

const (
	IOT_DEVICE_STATUS_END  = "iote-global-onoffline-end"
	IOT_DEVICE_STATUS_EDGE = "iote-global-onoffline-edge"
	DEVICE_STATUS_ONLINE   = "online"  // 在线
	DEVICE_STATUS_OFFLINE  = "offline" // 离线
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Infof("activity init setting:%+v", settings)
	options := initClientOption(ctx.Logger(), settings)

	if strings.HasPrefix(settings.Broker, "ssl") {

		cfg := &ssl.Config{}

		if len(settings.SSLConfig) != 0 {
			err := cfg.FromMap(settings.SSLConfig)
			if err != nil {
				return nil, err
			}

			if _, set := settings.SSLConfig["skipVerify"]; !set {
				cfg.SkipVerify = true
			}
			if _, set := settings.SSLConfig["useSystemCert"]; !set {
				cfg.UseSystemCert = true
			}
		} else {
			//using ssl but not configured, use defaults
			cfg.SkipVerify = true
			cfg.UseSystemCert = true
		}

		tlsConfig, err := ssl.NewClientTLSConfig(cfg)
		if err != nil {
			return nil, err
		}

		options.SetTLSConfig(tlsConfig)
	}

	mqttClient := mqtt.NewClient(options)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		ctx.Logger().Error(token.Error())
		return nil, token.Error()
	}

	act := &Activity{client: mqttClient, settings: settings}
	return act, nil
}

type Activity struct {
	settings *Settings
	client   mqtt.Client
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}

	err = ctx.GetInputObject(input)

	if err != nil {
		return true, err
	}
	params := make(map[string]interface{})
	if p, ok := input.Message["Params"]; ok {
		p, _ := coerce.ToObject(p)
		if v, ok := p["label"]; ok {
			params["id"] = "iotp-d09d0818-d957-49b6-9f96-0942fc3dd889"
			params["type"] = v.(map[string]interface{})["type"].(string)
			params["value"] = v.(map[string]interface{})["value"]
			params["time"] = time.Now().Unix() * 1000
			data := make(map[string]interface{})
			data["id"] = uuid.NewV4().String()
			data["version"] = "v1.0.0"
			data["params"] = params
			label := make(map[string]interface{})
			label["label"] = params
			data["params"] = label
			message, _ := json.Marshal(data)
			ctx.Logger().Infof("[Activity] Eval  Topic:%s,Message:%s", input.Topic, string(message))
			if token := a.client.Publish(a.settings.Topic, byte(a.settings.Qos), true, message); token.Wait() && token.Error() != nil {
				ctx.Logger().Debugf("Error in publishing: %v", err)
				return true, token.Error()
			}
		}
		//if v, ok := p["image"]; ok {
		//	params["id"] = v.(map[string]interface{})["id"].(string)
		//	params["type"] = v.(map[string]interface{})["type"].(string)
		//	params["value"] = v.(map[string]interface{})["value"]
		//	params["time"] = time.Now().Unix() * 1000
		//	data := make(map[string]interface{})
		//	data["id"] = uuid.NewV4().String()
		//	data["version"] = "v1.0.0"
		//
		//	image := make(map[string]interface{})
		//	image["image"] = params
		//	data["params"] = image
		//	message, _ := json.Marshal(data)
		//	ctx.Logger().Infof("[Activity] Eval  Topic:%s,Message:%s", input.Topic, string(message))
		//	if token := a.client.Publish(a.settings.Topic, byte(a.settings.Qos), true, message); token.Wait() && token.Error() != nil {
		//		ctx.Logger().Debugf("Error in publishing: %v", err)
		//		return true, token.Error()
		//	}
		//}
		if v, ok := p["color"]; ok {
			params["id"] = "iotp-bb73a69f-03a0-4305-9897-45a8ddf95e76"
			params["type"] = v.(map[string]interface{})["type"].(string)
			params["value"] = v.(map[string]interface{})["value"]
			params["time"] = time.Now().Unix() * 1000
			data := make(map[string]interface{})
			data["id"] = uuid.NewV4().String()
			data["version"] = "v1.0.0"
			color := make(map[string]interface{})
			color["color"] = params
			data["params"] = color

			message, _ := json.Marshal(data)
			ctx.Logger().Infof("[Activity] Eval  Topic:%s,Message:%s", input.Topic, string(message))
			if token := a.client.Publish(a.settings.Topic, byte(a.settings.Qos), true, message); token.Wait() && token.Error() != nil {
				ctx.Logger().Debugf("Error in publishing: %v", err)
				return true, token.Error()
			}
		}
		if v, ok := p["license"]; ok {
			params["id"] = "iotp-f9c651bc-aeb3-47b1-b393-484558377eb8"
			params["type"] = v.(map[string]interface{})["type"].(string)
			params["value"] = v.(map[string]interface{})["value"]
			params["time"] = time.Now().Unix() * 1000
			data := make(map[string]interface{})
			data["id"] = uuid.NewV4().String()
			data["version"] = "v1.0.0"
			license := make(map[string]interface{})
			license["license"] = params
			data["params"] = license
			message, _ := json.Marshal(data)
			ctx.Logger().Infof("[Activity] Eval  Topic:%s,Message:%s", input.Topic, string(message))
			if token := a.client.Publish(a.settings.Topic, byte(a.settings.Qos), true, message); token.Wait() && token.Error() != nil {
				ctx.Logger().Debugf("Error in publishing: %v", err)
				return true, token.Error()
			}
		}
	}
	//if token := a.client.Publish(input.Topic, byte(a.settings.Qos), true, input.Message); token.Wait() && token.Error() != nil {
	//	ctx.Logger().Debugf("Error in publishing: %v", err)
	//	return true, token.Error()
	//}
	ctx.Logger().Infof("Published Message Success: %+v", input.Message)

	return true, nil
}

func initClientOption(logger log.Logger, settings *Settings) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(settings.Broker)
	opts.SetClientID(settings.Id)
	opts.SetUsername(settings.Username)
	opts.SetPassword(settings.Password)
	opts.SetCleanSession(settings.CleanSession)

	if settings.Store != "" && settings.Store != ":memory:" {
		logger.Debugf("Using file store: %s", settings.Store)
		opts.SetStore(mqtt.NewFileStore(settings.Store))
	}

	return opts
}
