package encodex

import (
	"github.com/qingcloudhx/core/activity"
	uuid "github.com/satori/go.uuid"
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
func TestUUID(t *testing.T) {
	assert.NotNil(t, uuid.NewV4().String())
}
