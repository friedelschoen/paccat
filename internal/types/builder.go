package types

import (
	"strings"

	"friedelschoen.io/paccat/internal/ast"
)

type ValueBuilder struct {
	strings.Builder
	sources []StringSource
}

func (this *ValueBuilder) WriteValue(val *StringValue) {
	this.sources = append(this.sources, StringSource{Start: this.Len(), Len: len(val.Content), Value: val})
	this.WriteString(val.Content)
}

func (this *ValueBuilder) Value(source ast.Node) *StringValue {
	return &StringValue{
		Node:         source,
		Content:      this.String(),
		StringSource: this.sources,
	}
}
