package model

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
	settings := &Settings{}

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(settings, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	queryParams := map[string]string{
		"status": "ava",
	}
	tc.SetInput("queryParams", queryParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}
