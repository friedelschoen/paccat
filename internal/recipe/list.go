package recipe

import (
	"hash"
	"strings"
)

type recipeList struct {
	pos   Position
	items []Evaluable
}

func (this *recipeList) String() string {
	return "RecipeList"
}

func (this *recipeList) Eval(ctx Context) (Value, error) {
	return this, nil
}

func (this *recipeList) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	for _, value := range this.items {
		value.WriteHash(hash)
	}
}

func (this *recipeList) GetPosition() Position {
	return this.pos
}

func (this *recipeList) GetSource() Evaluable {
	return this
}

func (this *recipeList) GetName() string {
	return "list"
}

func (this *recipeList) ToString(ctx Context) (*StringValue, error) {
	builder := strings.Builder{}
	sources := []StringSource{}
	for _, item := range this.items {
		anyValue, err := item.Eval(ctx)
		if err != nil {
			return nil, WrapRecipeError(err, this.pos, "while evaluating list")
		}
		strValue, err := CastString(anyValue, ctx)
		if err != nil {
			return nil, WrapRecipeError(err, this.pos, "while evaluating list")
		}
		sources = append(sources, StringSource{builder.Len(), len(strValue.Content), strValue})
		builder.WriteString(strValue.Content)
	}
	return &StringValue{this, builder.String(), sources}, nil
}
