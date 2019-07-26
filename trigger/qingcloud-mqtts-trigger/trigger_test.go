package qingcloud_mqtts_trigger

import (
	"encoding/json"
	"testing"

	"github.com/qingcloudhx/core/action"
	"github.com/qingcloudhx/core/support"
	"github.com/qingcloudhx/core/support/test"
	"github.com/qingcloudhx/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "trigger-mqtt",
	"ref": "github.com/qingcloudhx/device-contrib/trigger/mqtt",
	"settings": {
		"autoreconnect": true,
        "broker": "ssl://mqtt.bosch-iot-hub.com:8883",
        "certstore": "*****8",
        "cleansess": false,
        "enableTLS": true,
				"id": "FlogoTest237123",
        "keepalive": 30,
        "password": "****",
        "store": ":memory:",
        "user": "*****"
    },
	"handlers": [
	  {
		"settings": {
			"topics": [
				"control/+/+/req/#"
			]
		},
		"action" : {
		  "id": "dummy"
		}
	  }
	]
  }`

func TestTrigger_Register(t *testing.T) {

	ref := support.GetRef(&Trigger{})
	f := trigger.GetFactory(ref)
	assert.NotNil(t, f)
}

func TestRestTrigger_Initialize(t *testing.T) {

	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
	})}

	trg, err := test.InitTrigger(f, config, actions)

	assert.Nil(t, err)
	assert.NotNil(t, trg)
	err = trg.Start()

	for {

	}
}
