package types

import (
	"fmt"
	"io"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

type Value interface {
	GetName() string
	GetSource() ast.Node
	ToJSON(Context, io.Writer) error
}

func CastValue[T Value](from Value) (T, error) {
	var empty T
	to, ok := from.(T)
	if !ok {
		return empty, errors.NewRecipeError(from.GetSource().GetPosition(), fmt.Sprintf("unable to convert %s to %s", from.GetName(), to.GetName()))
	}
	return to, nil
}

func CastString(from Value, ctx Context) (*StringValue, error) {
	strValue, ok := from.(StringLike)
	if !ok {
		return nil, errors.NewRecipeError(from.GetSource().GetPosition(), fmt.Sprintf("unable to convert %s to string", from.GetName()))
	}
	return strValue.ToString(ctx)
}

func CastBoolean(from Value, ctx Context) (bool, error) {
	boolValue, ok := from.(BooleanLike)
	if !ok {
		return false, errors.NewRecipeError(from.GetSource().GetPosition(), fmt.Sprintf("unable to convert %s to boolean", from.GetName()))
	}
	return boolValue.ToBoolean(ctx)
}

/* value-types */
type StringLike interface {
	ToString(Context) (*StringValue, error)
}

type BooleanLike interface {
	ToBoolean(Context) (bool, error)
}

type DictLike interface {
	GetAttrbute(string, Context) (Value, error)
}

type StringSource struct {
	Start int
	Len   int
	Value *StringValue /* underlying string-value */
}
