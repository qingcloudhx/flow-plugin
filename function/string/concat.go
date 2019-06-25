package string

import (
	"bytes"
	"fmt"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnConcat{})
}

type fnConcat struct {
}

func (fnConcat) Name() string {
	return "concat"
}

func (fnConcat) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, true
}

func (fnConcat) Eval(params ...interface{}) (interface{}, error) {
	if len(params) >= 2 {
		var buffer bytes.Buffer

		for _, v := range params {
			buffer.WriteString(v.(string))
		}
		return buffer.String(), nil
	}

	return "", fmt.Errorf("fnConcat function must have at least two arguments")
}
