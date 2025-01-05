package types

import (
	"fmt"
	"io"

	"friedelschoen.io/paccat/internal/ast"
)

type StringValue struct {
	source       ast.Node
	Content      string
	StringSource []StringSource
}

func (this *StringValue) GetSource() ast.Node {
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

func (this *StringValue) ToJSON(ctx Context, w io.Writer) error {
	fmt.Fprintf(w, "\"%s\"", this.Content)

	return nil
}
