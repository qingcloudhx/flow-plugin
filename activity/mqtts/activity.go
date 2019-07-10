package mqtts

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/muesli/cache2go"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/support/ssl"
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
	ctx.Logger().Infof("activity init")
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

	if settings.Topic == "" {
		if deviceId, thingId, err := parseToken(settings.Password); err != nil {
			return nil, err
		} else {
			settings.Topic = buildUpTopic(deviceId, thingId)
			ctx.Logger().Infof("mqtt topic:%s", settings.Topic)
		}
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
	if input.Type == "heartbeat" {
		data := &DeviceUpStatusMsg{}
		//decodeBytes, err := base64.StdEncoding.DecodeString(input.Message.(string))
		err = json.Unmarshal([]byte(input.Message.(string)), data)
		if err != nil {
			ctx.Logger().Errorf("Unmarshal error Message: %s", input.Message.(string))
			return false, nil
		}
		if _, err := get(data.DeviceId); err != nil {
			add(data.DeviceId, 15*time.Second, data, func(item *cache2go.CacheItem) {
				data.Status = DEVICE_STATUS_OFFLINE
				out, _ := json.Marshal(data)
				input.Message = out
				if token := a.client.Publish(input.Topic, byte(a.settings.Qos), true, input.Message); token.Wait() && token.Error() != nil {
					ctx.Logger().Debugf("Error in publishing: %v", err)
				} else {
					ctx.Logger().Debugf("Published Topic:%s, Message: %s", input.Topic, string(out))
				}
			})
		} else {
			ctx.Logger().Debugf("Recv Heartbeat Topic:%s,Message:%s", input.Topic, input.Message.(string))
			return true, nil
		}
		//add(data.DeviceId,15 * time.Second,data,func(key interface{}){})
	} else {

	}
	ctx.Logger().Infof("[Activity] Eval  Topic:%s,Message:%s", input.Topic, input.Message.(string))
	if token := a.client.Publish(input.Topic, byte(a.settings.Qos), true, input.Message); token.Wait() && token.Error() != nil {
		ctx.Logger().Debugf("Error in publishing: %v", err)
		return true, token.Error()
	}
	ctx.Logger().Debugf("Published Message Success: %s", input.Message.(string))

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
