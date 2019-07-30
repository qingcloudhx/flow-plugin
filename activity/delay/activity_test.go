package delay

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

	act := &Activity{Delay: 5000}
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{}
	err := tc.SetInputObject(input)
	assert.Nil(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

}
