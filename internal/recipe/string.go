package recipe

import (
	"fmt"
	"hash"
	"strings"
)

type recipeString struct {
	pos     position
	content []Evaluable
}

func (this *recipeString) String() string {
	builder := strings.Builder{}
	builder.WriteString("RecipeString{")
	for i, content := range this.content {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", content))
	}
	return builder.String()
}

func (this *recipeString) HasOutput() bool {
	for _, content := range this.content {
		if content.HasOutput() {
			return true
		}
	}
	return false
}

func (this *recipeString) Eval(ctx *Context, attr string) (string, error) {
	if attr != "" {
		return "", NoAttributeError{ctx, this.pos, "string", attr}
	}
	builder := strings.Builder{}
	for _, content := range this.content {
		str, err := content.Eval(ctx, "")
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}
	return builder.String(), nil
}

func (this *recipeString) WriteHash(hash hash.Hash) {
	hash.Write([]byte("string"))
	for _, content := range this.content {
		content.WriteHash(hash)
	}
}

func (this *recipeString) GetPosition() position {
	return this.pos
}
