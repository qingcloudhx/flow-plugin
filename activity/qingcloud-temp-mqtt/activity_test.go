package qingcloud_temp_mqtt

import (
	"testing"

	"github.com/qingcloudhx/core/activity"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}
