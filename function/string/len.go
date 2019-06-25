package string

import (
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnLen{})
}

type fnLen struct {
}

func (fnLen) Name() string {
	return "len"
}

func (fnLen) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnLen) Eval(params ...interface{}) (interface{}, error) {

	s := params[0].(string)

	return len(s), nil
}
