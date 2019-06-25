package string

import (
	"testing"

	"github.com/qingcloudhx/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFnEqualsIgnoreCase_Eval(t *testing.T) {
	f := &fnEqualsIgnoreCase{}

	v, err := function.Eval(f, "foo", "Bar")
	assert.Nil(t, err)
	assert.False(t, v.(bool))

	v, err = function.Eval(f, "foo", "Foo")
	assert.Nil(t, err)
	assert.True(t, v.(bool))
}
