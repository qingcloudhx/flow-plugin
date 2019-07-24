package expression

import (
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/resolve"
)

type Factory interface {
	NewExpr(exprStr string) (Expr, error)
}

type Expr interface {
	Eval(scope data.Scope) (interface{}, error)
}

type FactoryCreatorFunc func(resolve.CompositeResolver) Factory
