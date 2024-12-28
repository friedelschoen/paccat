package recipe

import (
	"fmt"
	"hash"
	"iter"
	"strings"
)

type recipeWith struct {
	pos          position
	dependencies Evaluable
	target       Evaluable
}

func (this *recipeWith) String() string {
	return fmt.Sprintf("RecipeWith{target=%v, depends=%v}", this.target, this.dependencies)
}

func indexSplit(original string, delimiter byte) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		start := 0
		for {
			// Find the next index of the delimiter
			index := strings.IndexByte(original[start:], delimiter)
			if index == -1 {
				// If no more delimiters, add the last part
				part := original[start:]
				if len(part) > 0 {
					yield(start, part)
				}
				break
			}

			// Get the actual index in the original string
			end := start + index
			part := original[start:end]
			if start != end {
				yield(end, part)
			}

			// Move the start index past the delimiter
			start = end + 1
		}
	}
}

func (this *recipeWith) Eval(ctx *Context, attr string) (string, []StringSource, error) {
	if attr != "" {
		return "", nil, NoAttributeError{ctx, this.pos, "with-statement", attr}
	}
	depend, source, err := this.dependencies.Eval(ctx, "")
	if err != nil {
		return "", nil, err
	}

	for indx, dep := range indexSplit(depend, ' ') {
		src := SourceAt(source, indx)[0]
		name := fmt.Sprintf("%016x", EvaluableSum(src))
		ctx.Database.Install(name, dep)
	}

	return this.target.Eval(ctx, "")
}

func (this *recipeWith) WriteHash(hash hash.Hash) {
	hash.Write([]byte("with"))
	this.dependencies.WriteHash(hash)
	this.target.WriteHash(hash)
}

func (this *recipeWith) GetPosition() position {
	return this.pos
}
