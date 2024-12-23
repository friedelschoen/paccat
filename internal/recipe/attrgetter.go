package recipe

import (
	"fmt"
	"hash"
)

type recipeGetter struct {
	pos       position
	target    Evaluable
	attribute string
}

func (this *recipeGetter) String() string {
	return fmt.Sprintf("RecipeGetter#%s{%v}", this.attribute, this.target)
}

func (this *recipeGetter) HasOutput() bool {
	return this.target.HasOutput()
}

func (this *recipeGetter) Eval(ctx *Context) (string, error) {
	value, err := this.target.Eval(ctx, this.attribute)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (this *recipeGetter) WriteHash(hash hash.Hash) {
	hash.Write([]byte("getter"))
	this.target.WriteHash(hash)
	hash.Write([]byte(this.attribute))
}

func (this *recipeGetter) GetPosition() position {
	return this.pos
}
