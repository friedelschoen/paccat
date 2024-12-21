package recipe

import (
	"fmt"
	"hash/maphash"
)

type RecipeStringLiteral struct {
	content string
}

func (this *RecipeStringLiteral) String() string {
	return fmt.Sprintf("RecipeStringLiteral#\"%s\"", string(this.content))
}

func (this *RecipeStringLiteral) HasOutput() bool {
	return false
}

func (this *RecipeStringLiteral) Build(ctx *Context) (string, error) {
	return string(this.content), nil
}

func (this *RecipeStringLiteral) WriteHash(hash maphash.Hash) {
	hash.WriteString(this.content)
}
