package ast

import (
	"fmt"
	"hash"
)

type recipeGetter struct {
	pos       Position
	target    Evaluable
	attribute string
}

func (this *recipeGetter) String() string {
	return fmt.Sprintf("RecipeGetter#%s{%v}", this.attribute, this.target)
}

func (this *recipeGetter) Eval(ctx Context) (Value, error) {
	anyValue, err := this.target.Eval(ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, fmt.Sprintf("while trying to get attribute `%s`", this.attribute))
	}

	dict, ok := anyValue.(DictLike)
	if !ok {
		return nil, NewRecipeError(anyValue.GetSource().GetPosition(), fmt.Sprintf("cannot cast %s to dict", anyValue.GetName()))
	}

	res, err := dict.GetAttrbute(this.attribute, ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, fmt.Sprintf("while trying to get attribute `%s`", this.attribute))
	}
	return res, nil
}

func (this *recipeGetter) WriteHash(hash hash.Hash) {
	hash.Write([]byte("getter"))
	this.target.WriteHash(hash)
	hash.Write([]byte(this.attribute))
}

func (this *recipeGetter) GetPosition() Position {
	return this.pos
}
