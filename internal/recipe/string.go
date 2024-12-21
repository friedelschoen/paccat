package recipe

import (
	"fmt"
	"hash"
	"strings"
)

type recipeString struct {
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

func (this *recipeString) Eval(ctx *Context) (string, error) {
	builder := strings.Builder{}
	for _, content := range this.content {
		str, err := content.Eval(ctx)
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}
	return builder.String(), nil
}

func (this *recipeString) WriteHash(hash hash.Hash) {
	for _, content := range this.content {
		content.WriteHash(hash)
	}
}
