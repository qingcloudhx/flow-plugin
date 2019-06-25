package number

import (
	"math/rand"
	"time"

	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnRandom{})
}

type fnRandom struct {
}

func (fnRandom) Name() string {
	return "random"
}

func (fnRandom) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeInt}, true
}

func (fnRandom) Eval(params ...interface{}) (interface{}, error) {

	limit := 10
	if len(params) > 0 {
		limit = params[0].(int)
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit), nil
}
