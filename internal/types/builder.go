package types

import (
	"strings"

	"friedelschoen.io/paccat/internal/ast"
)

type ValueBuilder struct {
	strings.Builder
	sources []StringSource
}

func (this *ValueBuilder) WriteValue(val *StringValue, quote bool) {
	content := val.Content
	if quote && strings.ContainsRune(content, ' ') {
		content = "\"" + strings.ReplaceAll(content, "\"", "\"\"") + "\""
	}
	this.sources = append(this.sources, StringSource{Start: this.Len(), Len: len(content), Value: val})
	this.WriteString(content)
}

func (this *ValueBuilder) Value(node ast.Node) *StringValue {
	return &StringValue{
		Node:         node,
		Content:      this.String(),
		StringSource: this.sources,
	}
}
