package ast

import (
	"hash"
	"strings"
)

type recipeList struct {
	pos   Position
	items []Evaluable
}

type ListValue recipeList

func (this *recipeList) String() string {
	return "RecipeList"
}

func (this *recipeList) Eval(ctx Context) (Value, error) {
	return (*ListValue)(this), nil
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

func (this *ListValue) GetSource() Evaluable {
	return (*recipeList)(this)
}

func (this *ListValue) GetName() string {
	return "list"
}

func (this *ListValue) ToString(ctx Context) (*StringValue, error) {
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
	return &StringValue{this.GetSource(), builder.String(), sources}, nil
}
