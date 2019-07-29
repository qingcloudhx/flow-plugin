package qingcloud_build_property

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
	s := make(map[string]interface{})
	s["xx"] = make(map[string]interface{})
	m := make(map[string]interface{})
	m["id"] = "1212"
	m["type"] = "int32"
	s["xx"] = m
	ss := Settings{s}
	act := &Activity{settings: &ss}
	tc := test.NewActivityContext(act.Metadata())

	input := &Input{}
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
