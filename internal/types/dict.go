package types

import (
	"fmt"
	"io"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

type DictValue ast.DictNode

func (this *DictValue) GetSource() ast.Node {
	return (*ast.DictNode)(this)
}

func (this *DictValue) GetName() string {
	return "dict"
}

func (this *DictValue) GetAttrbute(ctx Context, attr string) (Value, error) {
	eval, ok := this.Items[attr]
	if !ok {
		return nil, errors.NewRecipeError(this.GetSource().GetPosition(), fmt.Sprintf("unable to get `%s`, not defined in dict", attr))
	}
	value, err := ctx.Evaluate(eval.Value)
	if err != nil {
		return nil, errors.WrapRecipeError(err, this.GetSource().GetPosition(), fmt.Sprintf("while getting attribute `%s`", attr))
	}
	return value, nil
}

func (this *DictValue) ToJSON(ctx Context, w io.Writer) error {
	w.Write([]byte{'{'})
	i := 0
	for key, node := range this.Items {
		if i > 0 {
			w.Write([]byte{','})
		}
		fmt.Fprintf(w, "\"%s\":", key)
		value, err := ctx.Evaluate(node.Value)
		if err != nil {
			return err
		}
		value.ToJSON(ctx, w)
		i++
	}
	w.Write([]byte{'}'})

	return nil
}
