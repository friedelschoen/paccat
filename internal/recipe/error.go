package recipe

import (
	"errors"
	"fmt"
)

var (
	NoOutputError = errors.New("did not produce output")
)

func recipeError(ctx *Context, pos position, message string) string {
	return fmt.Sprintf("%s:%d:%d: %s", ctx.filename, pos.line, pos.col, message)
}

type NoAttributeError struct {
	ctx       *Context
	pos       position
	object    string
	attribute string
}

func (this NoAttributeError) Error() string {
	return recipeError(this.ctx, this.pos, fmt.Sprintf("%s cannot have attribute: %s", this.object, this.attribute))
}

type UnknownReferenceError struct {
	ctx  *Context
	pos  position
	name string
}

func (this UnknownReferenceError) Error() string {
	return recipeError(this.ctx, this.pos, fmt.Sprintf("unknown reference: %s", this.name))
}

type UnknownAttributeError struct {
	ctx  *Context
	pos  position
	file string
	name string
}

func (this UnknownAttributeError) Error() string {
	return recipeError(this.ctx, this.pos, fmt.Sprintf("recipe %s has no attribute: %s", this.file, this.name))
}
