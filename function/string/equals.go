package string

import (
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnEquals{})
}

type fnEquals struct {
}

func (fnEquals) Name() string {
	return "equals"
}

func (fnEquals) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnEquals) Eval(params ...interface{}) (interface{}, error) {

	s1 := params[0].(string)
	s2 := params[1].(string)
	return s1 == s2, nil
}
