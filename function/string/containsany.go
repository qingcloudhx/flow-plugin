package string

import (
	"strings"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	function.Register(&fnContainsAny{})
}

type fnContainsAny struct {
}

func (fnContainsAny) Name() string {
	return "containsAny"
}

func (fnContainsAny) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnContainsAny) Eval(params ...interface{}) (interface{}, error) {
	return strings.ContainsAny(params[0].(string), params[1].(string)), nil
}
