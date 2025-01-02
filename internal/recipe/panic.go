package recipe

import (
	"fmt"
	"hash"
)

type recipePanic struct {
	pos     Position
	message Evaluable
}

func (this *recipePanic) Eval(ctx Context) (Value, error) {
	value, err := this.message.Eval(ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while evaluating panic")
	}
	strValue, err := CastString(value, ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while evaluating panic")
	}

	return nil, NewRecipeError(this.pos, strValue.Content)
}

func (this *recipePanic) String() string {
	return fmt.Sprintf("RecipePanic")
}

func (this *recipePanic) WriteHash(hash hash.Hash) {
	hash.Write([]byte("panic"))
	this.message.WriteHash(hash)
}

func (this *recipePanic) GetPosition() Position {
	return this.pos
}
