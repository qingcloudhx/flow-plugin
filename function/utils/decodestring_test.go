package utils

import (
	"github.com/qingcloudhx/core/data/expression/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-6-26 上午9:41
 */
func TestFnDecodeString_Eval(t *testing.T) {
	f := fnDecodeString{}
	res, err := function.Eval(f, []byte("MTIxMjEy"))
	assert.Nil(t, err)
	t.Log(res)
}
