package ast

import (
	"fmt"
	"hash"
)

type recipeDict struct {
	pos   Position
	items map[string]Evaluable
}

type DictValue recipeDict

func (this *recipeDict) String() string {
	return "RecipeDict"
}

func (this *recipeDict) Eval(ctx Context) (Value, error) {
	return (*DictValue)(this), nil
}

func (this *recipeDict) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	for _, value := range this.items {
		value.WriteHash(hash)
	}
}

func (this *recipeDict) GetPosition() Position {
	return this.pos
}

func (this *DictValue) GetSource() Evaluable {
	return (*recipeDict)(this)
}

func (this *DictValue) GetName() string {
	return "dict"
}

func (this *DictValue) GetAttrbute(ctx Context, attr string) (Value, error) {
	eval, ok := this.items[attr]
	if !ok {
		return nil, NewRecipeError(this.pos, fmt.Sprintf("unable to get `%s`, not defined in dict", attr))
	}
	value, err := eval.Eval(ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, fmt.Sprintf("while getting attribute `%s`", attr))
	}
	return value, nil
}
