package mqtt

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
		"port":"1884",
		"tokenSsl":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyPuRZacRZoCgI8yQKf43\nHRiS08KM4d/q86b5VQ6Y0RGoRNPL3B8J9fxoZfXQO5I88Gnttp4S4OyD96iWiomI\nIQMr8DKmzZuRNxvj5GqYBdFtOfAFkwU/KdI2+QEJ2m37F6KZ+YiGu/+U6jMALwbK\naqHXS3mC+jlbYeK/FweZwc6pcEmtuLMn7rpBeogm6nyNF38oVBrh2zvYzqsovmc7\npLx4qZrk1GsJv8I/k/cPdH14DYgrKDGpjdbUIShoQllwrjR6mYRhpXlcr9OcimpC\nkgKJHif2WgMRBJfov9Ccg4eejN4iX0OUeIApA8YkyRre4dyEwfcWVDqUXsgVXIrt\n+QIDAQAB\n-----END PUBLIC KEY-----",
        "broker": "tcp://127.0.0.1:1883",
        "keepalive": 30,
		"id":"hexing"
    },
	"handlers": [
	  {
		"settings": {
			"topic_down":"down",
			"topic_up": "control/+/+/req/#"
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
		select {}
	}
}
