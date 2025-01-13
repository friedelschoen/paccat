package errors

import (
	"fmt"
	"io"
	"strings"
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
		err, ok := current.(Positioned)
		if !ok {
			fmt.Fprintf(writer, "??: %v\n", current)
		} else {
			pos := err.GetPosition()

			endOffset := 0
			line := 0
			var startLine, startOffset int

			lines := strings.SplitAfter(pos.File.Content, "\n")
			for _, lineStr := range lines {
				line++
				beginOffset := endOffset
				endOffset += len(lineStr)

				if pos.Start > endOffset {
					continue
				}

				if startLine == 0 {
					startLine = line
					startOffset = pos.Start - beginOffset
				}

				/* it's a oneliner */
				if pos.Start >= beginOffset && pos.End < endOffset {
					padding := 0
					for strings.ContainsRune(" \t", rune(lineStr[0])) {
						lineStr = lineStr[1:]
						padding--
					}
					fmt.Fprintf(writer, "%3d | %s", line, lineStr)
					if lineStr[len(lineStr)-1] != '\n' {
						fmt.Fprintln(writer)
					}
					writer.Write([]byte("    | ")) // Padding to align under text

					padding += pos.Start - beginOffset
					for i := 0; i < padding; i++ {
						writer.Write([]byte{' '})
					}
					writer.Write([]byte{'^'})

					length := pos.Len()
					for i := 0; i < length-1; i++ {
						writer.Write([]byte{'-'})
					}
					writer.Write([]byte{'\n'})
				} else {
					fmt.Fprintf(writer, "%3d |> %s", line, lineStr)
					if lineStr[len(lineStr)-1] != '\n' {
						fmt.Fprintln(writer)
					}
				}

				if pos.End < endOffset {
					break
				}
			}

			// Add the error message
			fmt.Fprintf(writer, "%s:%d:%d: %v\n", pos.File.Filename, startLine, startOffset+1, err)
		}
		prev, ok := current.(ContextError)
		if !ok {
			break
		}
		current = prev.Previous()
	}
}
