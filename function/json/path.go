package json

import (
	"github.com/oliveagle/jsonpath"
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnPath{})
}

type fnPath struct {
}

// Name returns the name of the function
func (fnPath) Name() string {
	return "path"
}

// Sig returns the function signature
func (fnPath) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeAny}, false
}

// Eval executes the function
func (fnPath) Eval(params ...interface{}) (interface{}, error) {
	expression := params[0].(string)
	return jsonpath.JsonPathLookup(params[1], expression)
}
