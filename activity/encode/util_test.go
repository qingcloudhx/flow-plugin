package encode

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
		Label:      "label",
		Confidence: 12.3,
		Image:      "sss",
		Color:      "blue",
		License:    "sas",
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
	res, err := buildMessage(in.ToMap(), settings)
	assert.Nil(t, err)
	t.Log(res)
}
