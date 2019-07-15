package qingcloud_channel

import (
	"testing"

	"github.com/qingcloudhx/core/activity"
	"github.com/qingcloudhx/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {

	act := &Activity{}
	tc := test.NewActivityContext(act.Metadata())

	input := &Input{Head: map[string]interface{}{"id": "test message"}, Body: []byte("assasas")}
	tc.SetInputObject(input)

	act.Eval(tc)
}

func TestAddToFlow(t *testing.T) {

	act := &Activity{}
	tc := test.NewActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("message", "test message")
	tc.SetInput("addDetails", true)

	act.Eval(tc)
}
