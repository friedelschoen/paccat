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

func (this *recipeGetter) Eval(ctx *Context) (string, []StringSource, error) {
	return this.target.Eval(ctx, this.attribute)
}

func (this *recipeGetter) WriteHash(hash hash.Hash) {
	hash.Write([]byte("getter"))
	this.target.WriteHash(hash)
	hash.Write([]byte(this.attribute))
}

func (this *recipeGetter) GetPosition() position {
	return this.pos
}
