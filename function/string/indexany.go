package string

import (
	"strings"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	function.Register(&fnIndexAny{})
}

type fnIndexAny struct {
}

func (fnIndexAny) Name() string {
	return "indexAny"
}

func (fnIndexAny) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (fnIndexAny) Eval(params ...interface{}) (interface{}, error) {
	return strings.IndexAny(params[0].(string), params[1].(string)), nil
}
