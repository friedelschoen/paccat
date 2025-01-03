package types

import (
	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

type LambdaValue ast.LambdaNode

func (this *LambdaValue) ToString(ctx Context) (*StringValue, error) {
	return nil, errors.NewRecipeError(this.GetSource().GetPosition(), "unable to convert a lambda to string")
}

func (this *LambdaValue) GetSource() ast.Node {
	return (*ast.LambdaNode)(this)
}

func (this *LambdaValue) GetName() string {
	return "lambda"
}
