package recipe

import (
	"fmt"
	"hash"
	"strings"
)

type recipeList struct {
	pos   position
	items []Evaluable
}

func (this *recipeList) String() string {
	builder := strings.Builder{}
	builder.WriteString("RecipeString{")
	for i, content := range this.items {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", content))
	}
	return builder.String()
}

func (this *recipeList) HasOutput() bool {
	for _, content := range this.items {
		if content.HasOutput() {
			return true
		}
	}
	return false
}

func (this *recipeList) Eval(ctx *Context, attr string) (string, error) {
	if attr != "" {
		return "", NoAttributeError{ctx, this.pos, "list", attr}
	}
	builder := strings.Builder{}
	for i, content := range this.items {
		if i > 0 {
			builder.WriteString(" ")
		}
		str, err := content.Eval(ctx, "")
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}
	return builder.String(), nil
}

func (this *recipeList) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	for _, value := range this.items {
		value.WriteHash(hash)
	}
}

func (this *recipeList) GetPosition() position {
	return this.pos
}
