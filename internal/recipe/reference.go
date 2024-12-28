package recipe

import (
	"fmt"
	"hash"
)

type recipeReference struct {
	pos  position
	name string
}

func (this *recipeReference) Eval(ctx *Context, attr string) (string, []StringSource, error) {
	return ctx.Get(this.name, attr)
}

func (this *recipeReference) String() string {
	return fmt.Sprintf("RecipeReference#%s", this.name)
}

func (this *recipeReference) WriteHash(hash hash.Hash) {
	hash.Write([]byte("reference"))
	hash.Write([]byte(this.name))
}

func (this *recipeReference) GetPosition() position {
	return this.pos
}
