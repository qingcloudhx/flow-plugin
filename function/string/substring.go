package string

import (
	"fmt"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnSubstring{})
}

type fnSubstring struct {
}

func (fnSubstring) Name() string {
	return "substring"
}

func (fnSubstring) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeInt, data.TypeInt}, false
}

func (fnSubstring) Eval(params ...interface{}) (interface{}, error) {

	str := params[0].(string)
	start := params[1].(int)
	length := params[2].(int)

	if length == -1 {
		return str[start:], nil
	}

	if start+length > len(str) {
		return nil, fmt.Errorf("string length exceeded")
	}

	return str[start : start+length], nil
}
