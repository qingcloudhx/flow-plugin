package qingcloud_encode_hearbeat

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/muesli/cache2go"
	"github.com/pkg/errors"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/support/ssl"
	"strings"
	"sync/atomic"
	"time"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var flag int32 = 0

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {
	ctx.Logger().Infof("activity init")
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		ctx.Logger().Error(err)
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
		go func() {
			for {
				if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
					ctx.Logger().Error(token.Error())
					time.Sleep(3 * time.Second)
				} else {
					ctx.Logger().Infof("mqtt conect success")
					atomic.StoreInt32(&flag, 1)
					break
				}
			}
		}()
	} else {
		ctx.Logger().Infof("mqtt conect success")
		atomic.StoreInt32(&flag, 1)
	}
	act := &Activity{settings: settings, client: mqttClient, log: ctx.Logger()} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	settings *Settings
	client   mqtt.Client
	log      log.Logger
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}
	ctx.Logger().Infof("Input DeviceId:%s,ThingId:%s", input.DeviceId, input.ThingId)

	if val := atomic.LoadInt32(&flag); val == 0 {
		return false, errors.New("mqtt client not init")
	}

	topic := fmt.Sprintf(a.settings.Format, input.ThingId, input.DeviceId)
	message := make(map[string]interface{})
	message["thingId"] = input.ThingId
	message["deviceID"] = input.DeviceId
	message["time"] = time.Now().Unix() * 1000
	message["status"] = "online"
	if exists(input.DeviceId) {
		return false, nil
	}
	add(topic, 15*time.Second, message, a.callback)
	data, _ := json.Marshal(message)
	if token := a.client.Publish(topic, byte(a.settings.Qos), true, data); token.Wait() && token.Error() != nil {
		ctx.Logger().Debugf("Error in publishing: %v", err)
		return false, token.Error()
	}

	output := &Output{Topic: topic, Message: string(data)}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
func (a *Activity) callback(item *cache2go.CacheItem) {
	message := item.Data()
	topic := ""
	if v, ok := message.(map[string]interface{}); ok {
		v["status"] = "offline"
		topic = fmt.Sprintf(a.settings.Format, v["thingId"], v["deviceId"])
		data, _ := json.Marshal(message)
		if token := a.client.Publish(topic, byte(a.settings.Qos), true, data); token.Wait() && token.Error() != nil {
			a.log.Errorf("Error in publishing: %v", token.Error())
		} else {
			a.log.Infof("[mqtt] thindId:%s,deviceId:%s offline", v["thingId"], v["deviceId"])
		}
	}
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
