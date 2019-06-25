package string

import (
	"strings"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	function.Register(&fnTrimLeft{})
}

type fnTrimLeft struct {
}

func (fnTrimLeft) Name() string {
	return "trimLeft"
}

func (fnTrimLeft) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnTrimLeft) Eval(params ...interface{}) (interface{}, error) {
	return strings.TrimLeft(params[0].(string), params[1].(string)), nil
}
