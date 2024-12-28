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

func (this *recipeList) Eval(ctx *Context, attr string) (string, []StringSource, error) {
	if attr != "" {
		return "", nil, NoAttributeError{ctx, this.pos, "list", attr}
	}
	builder := strings.Builder{}
	sources := []StringSource{}
	for i, content := range this.items {
		if i > 0 {
			builder.WriteString(" ")
		}
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

func (this *recipeList) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	for _, value := range this.items {
		value.WriteHash(hash)
	}
}

func (this *recipeList) GetPosition() position {
	return this.pos
}
