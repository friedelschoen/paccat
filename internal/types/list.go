package types

import (
	"strings"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

type ListValue ast.ListNode

func (this *ListValue) GetSource() ast.Node {
	return (*ast.ListNode)(this)
}

func (this *ListValue) GetName() string {
	return "list"
}

func (this *ListValue) ToString(ctx Context) (*StringValue, error) {
	builder := strings.Builder{}
	sources := []StringSource{}
	for _, item := range this.Items {
		anyValue, err := ctx.Evaluate(item)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetSource().GetPosition(), "while evaluating list")
		}
		strValue, err := CastString(anyValue, ctx)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetSource().GetPosition(), "while evaluating list")
		}
		sources = append(sources, StringSource{builder.Len(), len(strValue.Content), strValue})
		builder.WriteString(strValue.Content)
	}
	return &StringValue{this.GetSource(), builder.String(), sources}, nil
}
