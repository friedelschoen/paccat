package recipe

import (
	"fmt"
	"strings"
)

type RecipeStringLiteral struct {
	content string
}

type RecipeString struct {
	Elements []Buildable
	Previous string
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

func (this *RecipeString) String() string {
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

func (this *RecipeString) HasOutput() bool {
	for _, content := range this.Elements {
		if content.HasOutput() {
			return true
		}
	}
	return false
}

func (this *RecipeString) Build(ctx *Context) (string, error) {
	builder := strings.Builder{}
	for _, content := range this.Elements {
		str, err := content.Build(ctx)
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}
	return builder.String(), nil
}
