package modeldata

import (
	"fmt"
	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/data/mapper"
	"github.com/qingcloudhx/core/data/resolve"
	"github.com/qingcloudhx/core/support/test"
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
func TestAvtivity(t *testing.T) {
	settings := &Settings{DeviceId: "device-id", ThingId: "thing-id"}

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(settings, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	data := []DataModel{
		{
			Id:    "1212",
			Name:  "21212",
			Type:  "float",
			Value: 12.23,
		},
		{
			Id:    "1212",
			Name:  "21212",
			Type:  "int",
			Value: 121,
		},
	}

	tc.SetInput("device", data)

	//eval
	r, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, r)

	val := tc.GetOutput("message")
	fmt.Printf("result: %v\n", val)
}
