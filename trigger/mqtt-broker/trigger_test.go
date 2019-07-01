package mqtt_broker

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
        "broker": "tcp://192.168.14.120:8055",
        "keepalive": 30,
		"username":"hexing",
		"password":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItbDZncTBreGwiLCJlaXNrIjoiYjVYbVVReWxXNEVfUmZqSmp0bGlRU3lqcGF1S2dveDZ1eUlBbU9IbjNQQT0iLCJleHAiOjE1OTMyNTM4NjgsImlhdCI6MTU2MTcxNzg2OCwiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z1paSCIsIm5iZiI6MCwib3JnaSI6ImlvdGQtOGY4MzFjMzAtYjRjOS00YThhLTlmOGQtNTkzMmI5MDkwMmNlIiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoiaWFtX3JvbGUiLCJzdWIiOiJzdHMiLCJ0aGlkIjoiOCIsInR5cCI6IklEIn0.hWGMc_6hD3OugBgRLmEmnhro_pC5YXj7GGz0O8jSRZcnCepTDZfe9TJrjeHb4aZymrryqrhP5NU8R2Pa83kdbFJqAQY06avBC72Lu4Nhdx8bwE7inSCQo5os_xc7KBbjxgcoWYbeAGFAukAItw65-StMXl4G4E5kRfbXL0ioyfafLBYJM9IFmmVkPP50NDVsUhPhY8tiVzQhRn_P187ZKhvhOtkFH0T2iEeHmEx-vc3JVUIriHCB6eB5oAItsCd8YrIU4NnMhv1u42ecUDIYlWlFM29oSBfRg_5BEAbcb3pLFBQusIUfy3V4pMx3kYwfv5isiYOiJKXHcv-fag0hVQ",
		"id":"hexing",
		"autoReconnect":true
    },
	"handlers": [
	  {
		"settings": {
			"topic":"down"
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
	if err = trg.Start(); err != nil {
		return
	}
	for {
		select {}
	}
}
