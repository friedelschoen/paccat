package recipe

import (
	"fmt"
	"hash"
	"hash/crc64"
)

type Position struct {
	Filename string  /* name of file */
	Content  *string /* content of file */
	Start    int     /* begin-character */
	End      int     /* end of value */
}

func (this Position) Len() int {
	return this.End - this.Start
}

type StringSource struct {
	Start int
	Len   int
	Value *StringValue /* underlying string-value */
}

type Positioned interface {
	GetPosition() Position
}

type Named interface {
	GetName() string
}

type Evaluable interface {
	Positioned
	Eval(Context) (Value, error)
	WriteHash(hash.Hash)
}

type Value interface {
	Named
	GetSource() Evaluable
}

func EvaluableSum(in Evaluable) string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	in.WriteHash(hash)
	return fmt.Sprintf("%016x", hash.Sum64())
}

func CastValue[T Value](from Value) (T, error) {
	var empty T
	to, ok := from.(T)
	if !ok {
		return empty, NewRecipeError(from.GetSource().GetPosition(), fmt.Sprintf("unable to convert %s to %s", from.GetName(), to.GetName()))
	}
	return to, nil
}

func CastString(from Value, ctx Context) (*StringValue, error) {
	strValue, ok := from.(StringLike)
	if !ok {
		return nil, NewRecipeError(from.GetSource().GetPosition(), fmt.Sprintf("unable to convert %s to string", from.GetName()))
	}
	return strValue.ToString(ctx)
}

func CastBoolean(from Value, ctx Context) (bool, error) {
	boolValue, ok := from.(BooleanLike)
	if !ok {
		return false, NewRecipeError(from.GetSource().GetPosition(), fmt.Sprintf("unable to convert %s to boolean", from.GetName()))
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
