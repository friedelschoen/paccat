package recipe

import (
	"fmt"
	"hash/maphash"
)

type RecipeReference struct {
	content string
}

func (this *RecipeReference) Build(ctx *Context) (string, error) {
	return ctx.Get(this.content, false)
}

func (this *RecipeReference) HasOutput() bool {
	return this.content == "out"
}

func (this *RecipeReference) String() string {
	return fmt.Sprintf("RecipeReference#%s", this.content)
}

func (this *RecipeReference) WriteHash(hash maphash.Hash) {
	hash.WriteString(this.content)
}
