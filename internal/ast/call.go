package ast

import (
	"fmt"
	"hash"
)

type recipeCall struct {
	pos    Position
	target Evaluable
	args   map[string]Evaluable
}

func (this *recipeCall) String() string {
	return fmt.Sprintf("RecipeCall{%v}(%v)", this.target, this.args)
}

func (this *recipeCall) Eval(ctx Context) (Value, error) {
	value, err := this.target.Eval(ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while trying to call value")
	}
	lambda, err := CastValue[*LambdaValue](value)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while trying to call value")
	}

	newctx := ctx.Copy()
	for key, def := range lambda.args {
		if val, ok := this.args[key]; ok {
			newctx.scope[key] = val
		} else if def != nil {
			newctx.scope[key] = def
		} else {
			return nil, NewRecipeError(this.pos, fmt.Sprintf("lambda called without parameter `%s`", key))
		}
	}

	res, err := lambda.target.Eval(newctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while trying to call value")
	}
	return res, nil
}

func (this *recipeCall) WriteHash(hash hash.Hash) {
	hash.Write([]byte("call"))
	this.target.WriteHash(hash)
	for key, value := range this.args {
		hash.Write([]byte(key))
		value.WriteHash(hash)
	}
}

func (this *recipeCall) GetPosition() Position {
	return this.pos
}
