package recipe

import (
	"fmt"
	"hash"
)

type recipeWith struct {
	dependencies Evaluable
	target       Evaluable
}

func (this *recipeWith) String() string {
	return fmt.Sprintf("RecipeWith{target=%v, depends=%v}", this.target, this.dependencies)
}

func (this *recipeWith) HasOutput() bool {
	return this.dependencies.HasOutput() || this.target.HasOutput()
}

func (this *recipeWith) Eval(ctx *Context) (string, error) {
	return "", nil // TODO: ?
}

func (this *recipeWith) WriteHash(hash hash.Hash) {
	this.dependencies.WriteHash(hash)
	this.target.WriteHash(hash)
}
