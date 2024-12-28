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
	return fmt.Sprintf("RecipeStringLiteral#\"%s\"", this.value)
}

func (this *recipeStringLiteral) Eval(ctx *Context, attr string) (string, []StringSource, error) {
	if attr != "" {
		return "", nil, NoAttributeError{ctx, this.pos, "literal", attr}
	}
	return this.value, []StringSource{{0, len(this.value), this}}, nil
}

func (this *recipeStringLiteral) WriteHash(hash hash.Hash) {
	hash.Write([]byte("literal"))
	hash.Write([]byte(this.value))
}

func (this *recipeStringLiteral) GetPosition() position {
	return this.pos
}
