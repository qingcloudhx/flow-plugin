package encodex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-7-17 上午11:03
 */
func TestBuild(t *testing.T) {
	in := Input{
		Message: map[string]interface{}{
			"label":      "label",
			"confidence": 12.3,
			"image":      "sss",
		},
	}
	settings := map[string]interface{}{
		"label": map[string]interface{}{
			"id":   "1",
			"type": "string",
		},
		"confidence": map[string]interface{}{
			"id":   "2",
			"type": "float",
		},
		"image": map[string]interface{}{
			"id":   "3",
			"type": "string",
		},
	}
	res, err := buildMessage(in.Message, settings)
	assert.Nil(t, err)
	t.Log(res)
}
