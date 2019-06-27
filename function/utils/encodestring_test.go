package utils

import (
	"github.com/qingcloudhx/core/data/expression/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-6-26 上午9:40
 */
func TestFnEncodeString_Eval(t *testing.T) {
	f := &fnEncodeString{}
	res, err := function.Eval(f, []byte("121212"))
	assert.Nil(t, err)
	t.Log(res)
}
