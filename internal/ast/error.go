package ast

import (
	"fmt"
	"io"
	"strings"
)

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

func NewRecipeError(pos Position, message string) error {
	return &RecipeError{pos, nil, message}
}

func WrapRecipeError(previous error, pos Position, message string) error {
	if _, ok := previous.(*RecipeError); ok {
		return &RecipeError{pos, previous, message}
	} else if message != "" {
		return &RecipeError{pos, nil, fmt.Sprintf("%s: %v", message, previous)}
	} else {
		return &RecipeError{pos, nil, previous.Error()}
	}
}

func PrintTrace(writer io.Writer, current error) {
	for current != nil {
		err, ok := current.(*RecipeError)
		if !ok {
			fmt.Fprintf(writer, "??: %v\n", current)
			break
		}

		endOffset := 0
		line := 0
		var startLine, startOffset int

		lines := strings.SplitAfter(*err.pos.Content, "\n")
		for _, lineStr := range lines {
			line++
			beginOffset := endOffset
			endOffset += len(lineStr)

			if err.pos.Start > endOffset {
				continue
			}

			if startLine == 0 {
				startLine = line
				startOffset = err.pos.Start - beginOffset
			}

			/* it's a oneliner */
			if err.pos.Start >= beginOffset && err.pos.End < endOffset {
				fmt.Fprintf(writer, "%3d | %s", line, lineStr)
				writer.Write([]byte("    | ")) // Padding to align under text

				padding := err.pos.Start - beginOffset
				length := err.pos.Len()

				for i := 0; i < padding; i++ {
					writer.Write([]byte{' '})
				}
				writer.Write([]byte{'^'})
				for i := 0; i < length-1; i++ {
					writer.Write([]byte{'-'})
				}
				writer.Write([]byte{'\n'})
			} else {
				fmt.Fprintf(writer, "%3d |> %s", line, lineStr)
			}

			if err.pos.End < endOffset {
				break
			}
		}

		// Add the error message
		fmt.Fprintf(writer, "%s:%d:%d: %s\n", err.pos.Filename, startLine, startOffset+1, err.message)
		current = err.previous
	}
}
