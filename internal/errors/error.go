package errors

import (
	"fmt"
)

type ErrorFile struct {
	Filename string /* name of file */
	Content  string /* content of file */
}

type Position struct {
	File  *ErrorFile
	Start int /* begin-character */
	End   int /* end of value */
}

func (this Position) Len() int {
	return this.End - this.Start
}

func (this Position) Stretch(to Position) Position {
	return Position{
		File:  this.File,
		Start: this.Start,
		End:   to.End,
	}
}

type Positioned interface {
	GetPosition() Position
}

type ContextError interface {
	Previous() error
}

type RecipeError struct {
	pos      Position
	previous error
	message  string
}

func (this *RecipeError) Error() string {
	return this.message
}

func (this *RecipeError) GetPosition() Position {
	return this.pos
}

func (this *RecipeError) Previous() error {
	return this.previous
}

func NewRecipeError(pos Position, message string) error {
	return &RecipeError{pos, nil, message}
}

func WrapRecipeError(previous error, pos Position, message string) error {
	switch previous.(type) {
	case ContextError, Positioned:
		return &RecipeError{pos, previous, message}
	default:
		if message != "" {
			return &RecipeError{pos, nil, fmt.Sprintf("%s: %v", message, previous)}
		} else {
			return &RecipeError{pos, nil, previous.Error()}
		}
	}
}
