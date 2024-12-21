package recipe

import (
	"fmt"
	"hash/maphash"
	"strings"
)

type RecipeList struct {
	Elements []Buildable
}

func (this *RecipeList) String() string {
	builder := strings.Builder{}
	builder.WriteString("RecipeString{")
	for i, content := range this.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", content))
	}
	return builder.String()
}

func (this *RecipeList) HasOutput() bool {
	for _, content := range this.Elements {
		if content.HasOutput() {
			return true
		}
	}
	return false
}

func (this *RecipeList) Build(ctx *Context) (string, error) {
	builder := strings.Builder{}
	for i, content := range this.Elements {
		if i > 0 {
			builder.WriteString(" ")
		}
		str, err := content.Build(ctx)
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}
	return builder.String(), nil
}

func (this *RecipeList) WriteHash(hash maphash.Hash) {
	for _, value := range this.Elements {
		value.WriteHash(hash)
	}
}
