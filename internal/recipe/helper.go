package recipe

//go:generate pigeon -o parser.go recipe.peg

import (
	"cmp"
	"hash"
	"log"
	"maps"
	"reflect"
	"slices"
	"strings"
)

type pair struct {
	key   string
	value Evaluable
}

func asString(val any) string {
	builder := strings.Builder{}
	for _, chars := range val.([]any) {
		builder.Write(chars.([]byte))
	}
	return builder.String()
}

func makeString(val any) *recipeString {
	builder := strings.Builder{}
	result := make([]Evaluable, 0)
	for _, content := range val.([]any) {
		switch element := content.(type) {
		case []byte:
			builder.Write(element)
		case Evaluable:
			if builder.Len() > 0 {
				result = append(result, &recipeStringLiteral{builder.String()}, element)
				builder.Reset()
			} else {
				result = append(result, element)
			}
		default:
			log.Panicf("unexpected element: %v\n", reflect.TypeOf(element))
		}
	}
	if builder.Len() > 0 {
		result = append(result, &recipeStringLiteral{builder.String()})
	}
	return &recipeString{result}
}

// Combine head and tail into a single slice
func headTail[T any](head any, tail []any) []T {
	result := make([]T, 1+len(tail))
	result[0] = head.(T)
	for i, t := range tail {
		result[i+1] = t.(T)
	}
	return result
}

func toAnySlice[T any](input []any) []T {
	result := make([]T, len(input))
	for i, e := range input {
		if e != nil {
			result[i] = e.(T)
		}
	}
	return result
}

func collectPairs(pairs []pair) map[string]Evaluable {
	result := make(map[string]Evaluable, len(pairs))
	for _, keyvalue := range pairs {
		result[keyvalue.key] = keyvalue.value
	}
	return result
}

func sortedKeys[key cmp.Ordered, value any](attr map[key]value) []key {
	keys := make([]key, 0, len(attr))
	slices.AppendSeq(keys, maps.Keys(attr))
	slices.Sort(keys)
	return keys
}

func writeHashMap[key cmp.Ordered](attr map[key]Evaluable, hash hash.Hash) {
	keys := make([]key, 0, len(attr))
	keys = slices.AppendSeq(keys, maps.Keys(attr))
	slices.Sort(keys)

	for _, key := range keys {
		if attr[key] != nil {
			attr[key].WriteHash(hash)
		}
	}
}
