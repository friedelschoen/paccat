package recipe

import "strings"

type Wrappable interface {
	Wrap() (Buildable, error)
}

type recipe struct {
	require, values []pair
}

type pair struct {
	key   key
	value Wrappable
}

type functionCall struct {
	path   Wrappable
	params []pair
}

type attrGetter struct {
	call Wrappable
	attr key
}

type key string
type recipeString []any
type recipeLiteral string

func asString(val any) string {
	builder := strings.Builder{}
	for _, chars := range val.([]any) {
		builder.Write(chars.([]byte))
	}
	return builder.String()
}

func makeString(val any) recipeString {
	builder := strings.Builder{}
	result := make([]any, 0)
	for _, content := range val.([]any) {
		if chars, ok := content.([]byte); ok {
			builder.Write(chars)
		} else if builder.Len() > 0 {
			result = append(result, recipeLiteral(builder.String()), content)
			builder.Reset()
		} else {
			result = append(result, content)
		}
	}
	if builder.Len() > 0 {
		result = append(result, recipeLiteral(builder.String()))
	}
	return recipeString(result)
}

// Combine head and tail into a single slice
func headTail[T any](head any, tail []any) []T {
	result := make([]T, 0, 1+len(tail))
	result = append(result, head.(T))
	for _, t := range tail {
		result = append(result, t.(T))
	}
	return result
}

func toAnySlice[T any](input []any) []T {
	result := make([]T, len(input))
	for i, e := range input {
		result[i] = e.(T)
	}
	return result
}
