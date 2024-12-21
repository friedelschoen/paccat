package recipe

import (
	"fmt"
	"hash"
	"strings"
)

type recipeList struct {
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

func (this *recipeList) Eval(ctx *Context) (string, error) {
	builder := strings.Builder{}
	for i, content := range this.items {
		if i > 0 {
			builder.WriteString(" ")
		}
		str, err := content.Eval(ctx)
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}
	return builder.String(), nil
}

func (this *recipeList) WriteHash(hash hash.Hash) {
	for _, value := range this.items {
		value.WriteHash(hash)
	}
}
