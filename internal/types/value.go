package types

import (
	"friedelschoen.io/paccat/internal/ast"
)

type StringValue struct {
	Node         ast.Node
	Content      string
	StringSource []StringSource
	Attributes   map[string]ValuePair
}

type ValuePair struct {
	Key   *StringValue
	Value *StringValue
}

type StringSource struct {
	Start int
	Len   int
	Value *StringValue /* underlying string-value */
}

func (this *StringValue) ValueAt(pos int) *StringValue {
	for _, item := range this.StringSource {
		if pos >= item.Start && pos < item.Start+item.Len {
			return item.Value
		}
	}
	return this
}
