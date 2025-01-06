package types

import (
	"strings"

	"friedelschoen.io/paccat/internal/ast"
)

type ValueBuilder struct {
	strings.Builder

	sources []StringSource
}

func (this *ValueBuilder) WriteValue(val *StringValue, quoted bool) {
	str := val.Content
	if quoted {
		str = val.Quoted()
	}

	this.sources = append(this.sources, StringSource{Start: this.Len(), Len: len(str), Value: val})
	this.WriteString(str)
}

func (this *ValueBuilder) Value(source ast.Node) *StringValue {
	return &StringValue{
		source:       source,
		Content:      this.String(),
		StringSource: this.sources,
	}
}
