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

func (this *recipeString) Eval(ctx *Context, attr string) (string, []StringSource, error) {
	if attr != "" {
		return "", nil, NoAttributeError{ctx, this.pos, "string", attr}
	}
	builder := strings.Builder{}
	sources := []StringSource{}
	for _, content := range this.content {
		str, strSources, err := content.Eval(ctx, "")
		if err != nil {
			return "", nil, err
		}
		sources = append(sources, offsetSources(strSources, builder.Len())...)
		builder.WriteString(str)
	}
	sources = append(sources, StringSource{0, builder.Len(), this})
	return builder.String(), sources, nil
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
