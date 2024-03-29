package qingcloud_mqtt_client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/support/ssl"
	"strings"
	"sync/atomic"
	"time"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

var flag int32 = 0

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

	if token := mqttClient.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		ctx.Logger().Error(token.Error())
		go func() {
			for {
				if token := mqttClient.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
					ctx.Logger().Error(token.Error())
					time.Sleep(3 * time.Second)
				} else {
					ctx.Logger().Infof("mqtt conect success")
					atomic.StoreInt32(&flag, 1)
					break
				}
			}
		}()
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

	if token := a.client.Publish(a.settings.Topic, byte(a.settings.Qos), true, input.Message); token.Wait() && token.Error() != nil {
		ctx.Logger().Debugf("Error in publishing: %v", err)
		return true, token.Error()
	}

	ctx.Logger().Infof("Published Topic:%s,Message: %v", a.settings.Topic, input.Message)

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
