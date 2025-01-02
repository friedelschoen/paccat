package recipe

import (
	"fmt"
	"hash"
	"strings"
)

type recipeString struct {
	pos     Position
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

func (this *recipeString) Eval(ctx Context) (Value, error) {
	builder := strings.Builder{}
	sources := []StringSource{}
	for _, content := range this.content {
		value, err := content.Eval(ctx)
		if err != nil {
			return nil, err
		}
		strValue, err := CastString(value, ctx)
		if err != nil {
			return nil, err
		}
		sources = append(sources, StringSource{builder.Len(), len(strValue.Content), strValue})
		builder.WriteString(strValue.Content)
	}
	return &StringValue{this, builder.String(), sources}, nil
}

func (this *recipeString) WriteHash(hash hash.Hash) {
	hash.Write([]byte("string"))
	for _, content := range this.content {
		content.WriteHash(hash)
	}
}

func (this *recipeString) GetPosition() Position {
	return this.pos
}

type StringValue struct {
	source       Evaluable
	Content      string
	StringSource []StringSource
}

func (this *StringValue) GetSource() Evaluable {
	return this.source
}

func (this *StringValue) GetName() string {
	return "string"
}

func (this *StringValue) ToString(ctx Context) (*StringValue, error) {
	return this, nil
}

func (this *StringValue) ValueAt(pos int) *StringValue {
	for _, item := range this.StringSource {
		if pos >= item.Start && pos < item.Start+item.Len {
			return item.Value
		}
	}
	return this
}
