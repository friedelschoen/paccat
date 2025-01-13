package types

import (
	"iter"
	"strings"

	"friedelschoen.io/paccat/internal/ast"
)

type StringValue struct {
	Node         ast.Node
	Content      string
	StringSource []StringSource
	Attributes   map[string]ValuePair
}

type ValuePair struct {
	Key   *StringValue
	Value *StringValue
}

type StringSource struct {
	Start int
	Len   int
	Value *StringValue /* underlying string-value */
}

func (this *StringValue) ValueAt(pos int) *StringValue {
	for _, item := range this.StringSource {
		if pos >= item.Start && pos < item.Start+item.Len {
			return item.Value
		}
	}
	return this
}

func (this *StringValue) Split() iter.Seq2[string, *StringValue] {
	return func(yield func(string, *StringValue) bool) {
		begin := 0
		inString := false
		builder := strings.Builder{}

		for i := 0; i <= len(this.Content); i++ {
			// Handle end of content for final yield
			if i == len(this.Content) || (!inString && strings.ContainsRune(" \t\n\r", rune(this.Content[i]))) {
				if i > begin { // Ensure non-empty substring
					var value *StringValue = nil
					for _, item := range this.StringSource {
						if item.Start <= begin && item.Start+item.Len >= i {
							value = item.Value
							break
						}
					}

					// Handle string from builder or direct slice
					part := builder.String()
					if part == "" {
						part = this.Content[begin:i]
					} else {
						builder.Reset()
					}

					if !yield(part, value) {
						return
					}
				}
				begin = i + 1
			} else if inString {
				if this.Content[i] == '"' {
					if i < len(this.Content)-1 && this.Content[i+1] == '"' {
						// Escaped quote
						builder.WriteByte('"')
						i++
					} else {
						// End of quoted string
						inString = false
					}
				} else {
					builder.WriteByte(this.Content[i])
				}
			} else {
				if this.Content[i] == '"' {
					// Start of quoted string
					inString = true
					begin = i + 1
				}
			}
		}
	}
}

func (this *StringValue) FlatSources() iter.Seq[StringSource] {
	return func(yield func(StringSource) bool) {
		if !yield(StringSource{0, len(this.Content), this}) {
			return
		}
		for _, source := range this.StringSource {
			if !yield(source) {
				return
			}
			for child := range source.Value.FlatSources() {
				if !yield(StringSource{source.Start + child.Start, child.Len, child.Value}) {
					return
				}
			}
		}
	}
}
