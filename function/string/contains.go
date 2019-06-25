package string

import (
	"strings"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnContains{})
}

type fnContains struct {
}

func (s *fnContains) Name() string {
	return "contains"
}

func (fnContains) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnContains) Eval(params ...interface{}) (interface{}, error) {
	str1 := params[0].(string)
	str2 := params[1].(string)
	return strings.Contains(str1, str2), nil
}
