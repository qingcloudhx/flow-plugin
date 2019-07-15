package qingcloud_mqtt

import (
	"encoding/json"
	"github.com/qingcloudhx/core/action"
	"github.com/qingcloudhx/core/support/test"
	"github.com/qingcloudhx/core/trigger"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-6-27 上午10:54
 */
const testConfig string = `{
	"id": "flogo-model-trigger",
	"ref": "qingcloud-flow/plugin/trigger/mqttbroker",
	"settings": {
		"url":"tcp:127.0.0.1:1083",
		"event":"hehehe"
	},
	"handlers": [
	  {
			"action":{
				"id":"dummy"
			},
			"settings": {
		  	
			}
	  }
	]
	
  }`

func TestFactory_New(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	err = trg.Start()
	assert.Nil(t, err)
	err = trg.Stop()
	assert.Nil(t, err)
}
