package encode

import (
	"github.com/qingcloudhx/core/activity"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-7-4 下午5:48
 */
func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}
