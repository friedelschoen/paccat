package recipe

import (
	"fmt"
	"hash"
)

type recipeGetter struct {
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
	ctx.importAttribute = &this.attribute
	value, err := this.target.Eval(ctx)
	if err != nil {
		return "", err
	}
	if ctx.importAttribute != nil {
		return "", fmt.Errorf("attribute-getter not applied on function")
	}
	return value, nil
}

func (this *recipeGetter) WriteHash(hash hash.Hash) {
	this.target.WriteHash(hash)
	hash.Write([]byte(this.attribute))
}
