package recipe

import (
	"fmt"
	"hash"
)

type recipeStringLiteral struct {
	pos   position
	value string
}

func (this *recipeStringLiteral) String() string {
	return fmt.Sprintf("RecipeStringLiteral#\"%s\"", string(this.value))
}

func (this *recipeStringLiteral) Eval(ctx *Context, attr string) (string, error) {
	if attr != "" {
		return "", NoAttributeError{ctx, this.pos, "literal", attr}
	}
	return string(this.value), nil
}

func (this *recipeStringLiteral) WriteHash(hash hash.Hash) {
	hash.Write([]byte("literal"))
	hash.Write([]byte(this.value))
}

func (this *recipeStringLiteral) GetPosition() position {
	return this.pos
}
