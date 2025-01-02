package recipe

import (
	"fmt"
	"hash"
)

type recipeLambda struct {
	pos    Position
	target Evaluable
	args   map[string]Evaluable
}

func (this *recipeLambda) Eval(ctx Context) (Value, error) {
	return this, nil
}

func (this *recipeLambda) String() string {
	return fmt.Sprintf("RecipeLambda#{%v}{%v}", this.args, this.target)
}

func (this *recipeLambda) WriteHash(hash hash.Hash) {
	hash.Write([]byte("lambda"))
	this.target.WriteHash(hash)
	for key, value := range this.args {
		hash.Write([]byte(key))
		if value != nil {
			value.WriteHash(hash)
		}
	}
}

func (this *recipeLambda) GetPosition() Position {
	return this.pos
}

func (this *recipeLambda) ToString(ctx Context) (*StringValue, error) {
	return nil, NewRecipeError(this.pos, "unable to convert a lambda to string")
}

func (this *recipeLambda) GetSource() Evaluable {
	return this
}

func (this *recipeLambda) GetName() string {
	return "lambda"
}
