package qingcloud_tcp

import (
	"encoding/json"
	"testing"

	"github.com/qingcloudhx/core/action"
	"github.com/qingcloudhx/core/support/test"
	"github.com/qingcloudhx/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "flogo-tcpudp-trigger",
	"ref": "qingcloud-flow/plugin/trigger/tcp",
	"settings": {
      "network": "tcp",
	  "host": "127.0.0.1",
      "port": "8982"
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

func TestTrigger_Initialize(t *testing.T) {
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
