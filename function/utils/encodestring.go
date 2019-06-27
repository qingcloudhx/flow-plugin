package utils

import (
	"encoding/base64"
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

const (
	param_length = 1
)

func init() {
	function.Register(&fnEncodeString{})
}

type fnEncodeString struct {
}

func (fnEncodeString) Name() string {
	return "encodeString"
}

func (fnEncodeString) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBytes}, false
}

// Eval - UUID generates a random UUID according to RFC 4122
func (fnEncodeString) Eval(params ...interface{}) (interface{}, error) {
	str := base64.StdEncoding.EncodeToString(params[0].([]byte))
	return str, nil
}
