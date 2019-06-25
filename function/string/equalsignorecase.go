package string

import (
	"strings"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnEqualsIgnoreCase{})
}

type fnEqualsIgnoreCase struct {
}

func (s *fnEqualsIgnoreCase) Name() string {
	return "equalsIgnoreCase"
}

func (fnEqualsIgnoreCase) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnEqualsIgnoreCase) Eval(params ...interface{}) (interface{}, error) {
	str1 := params[0].(string)
	str2 := params[1].(string)
	return strings.EqualFold(str1, str2), nil
}
