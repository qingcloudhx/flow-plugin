package qingcloud_temp_mqtt

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

func init() {
	_ = activity.Register(&Activity{}, New)
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

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
		return nil, token.Error()
	}

	//if settings.Topic == "" {
	//	if deviceId, thingId, err := parseToken(settings.Password); err != nil {
	//		return nil, err
	//	} else {
	//		settings.Topic = buildUpTopic(deviceId, thingId)
	//		ctx.Logger().Infof("mqtt topic:%s", settings.Topic)
	//	}
	//}
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
	param, err := coerce.ToObject(input.Message)
	if err != nil {
		return true, err
	}
	ctx.Logger().Infof("eval start:%+v", param)
	if p, ok := param["params"]; ok {
		params := make(map[string]interface{})
		ps, _ := coerce.ToObject(p)
		if v, ok := ps["dt"]; ok {
			if val, ok := v.(map[string]interface{}); ok {
				val["id"] = "78"
				params["type"] = "string"
				if input.Color == "blue" {
					params["value"] = "距离<=20cm"
				} else {
					params["value"] = "距离>20cm"
				}
			}
		}
		params["id"] = "iotp-32f168bc-db95-4eef-851b-28cf1fd75712"
		params["time"] = time.Now().Unix() * 1000
		data := make(map[string]interface{})
		data["id"] = uuid.NewV4().String()
		data["version"] = "v1.0.0"
		dt := make(map[string]interface{})
		dt["dt"] = params
		data["params"] = dt
		message, _ := json.Marshal(data)
		ctx.Logger().Infof("eval topic:/sys/iott-bbeebd96-328e-4076-a59e-5a8341f5ab88/iotd-f6f1627e-ab18-49af-9d1c-88062ba44390/thing/event/property/post format:%+v", string(message))
		if token := a.client.Publish("/sys/iott-bbeebd96-328e-4076-a59e-5a8341f5ab88/iotd-f6f1627e-ab18-49af-9d1c-88062ba44390/thing/event/property/post", byte(a.settings.Qos), true, message); token.Wait() && token.Error() != nil {
			ctx.Logger().Debugf("Error in publishing: %v", err)
			return true, token.Error()
		}
		if v, ok := ps["temperature"]; ok {
			if val, ok := v.(map[string]interface{}); ok {
				val["id"] = "79"
			}
		}
	}

	color := make(map[string]interface{})
	color["id"] = "77"
	color["type"] = "string"
	color["value"] = input.Color
	color["time"] = time.Now().Unix() * 1000
	param["color"] = color

	power := make(map[string]interface{})
	power["id"] = "76"
	power["type"] = "float"
	power["value"] = 1
	power["time"] = time.Now().Unix() * 1000
	param["power"] = power
	message, _ := json.Marshal(param)
	ctx.Logger().Infof("eval event topic:%s,format:%+v", a.settings.Topic, string(message))
	if token := a.client.Publish(a.settings.Topic, byte(a.settings.Qos), true, message); token.Wait() && token.Error() != nil {
		ctx.Logger().Debugf("Error in publishing: %v", err)
		return true, token.Error()
	}

	ctx.Logger().Debugf("Published Message: %v", input.Message)

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
