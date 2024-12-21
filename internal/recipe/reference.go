package recipe

import (
	"fmt"
	"hash"
)

type recipeReference struct {
	name string
}

func (this *recipeReference) Eval(ctx *Context) (string, error) {
	return ctx.Get(this.name, false)
}

func (this *recipeReference) HasOutput() bool {
	return this.name == "out"
}

func (this *recipeReference) String() string {
	return fmt.Sprintf("RecipeReference#%s", this.name)
}

func (this *recipeReference) WriteHash(hash hash.Hash) {
	hash.Write([]byte(this.name))

}
