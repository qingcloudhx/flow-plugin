package expression

import (
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/resolve"
)

func NewLiteralExpr(val interface{}) Expr {
	return &literalExpr{val: val}
}

type literalExpr struct {
	val interface{}
}

func (e *literalExpr) Eval(scope data.Scope) (interface{}, error) {
	return e.val, nil
}

type resolutionExpr struct {
	resolution resolve.Resolution
}

func (e *resolutionExpr) Eval(scope data.Scope) (interface{}, error) {

	return e.resolution.GetValue(scope)
}
