package recipe

import (
	"fmt"
	"hash"
	"strings"
)

type recipeWith struct {
	dependencies Evaluable
	target       Evaluable
}

func (this *recipeWith) String() string {
	return fmt.Sprintf("RecipeWith{target=%v, depends=%v}", this.target, this.dependencies)
}

func (this *recipeWith) HasOutput() bool {
	return true // this must always build
}

func (this *recipeWith) Eval(ctx *Context) (string, error) {
	depends, err := this.dependencies.Eval(ctx)
	if err != nil {
		return "", err
	}

	for _, dep := range strings.Split(depends, " ") {
		err := installPath(dep)
		if err != nil {
			return "", err
		}
	}

	return this.target.Eval(ctx)
}

func (this *recipeWith) WriteHash(hash hash.Hash) {
	this.dependencies.WriteHash(hash)
	this.target.WriteHash(hash)
}
