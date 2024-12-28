package recipe

import (
	"fmt"
	"hash"
	"hash/crc64"
)

type StringSource struct {
	Start int
	Len   int
	Value Evaluable
}

type Evaluable interface {
	Eval(*Context, string) (string, []StringSource, error)
	WriteHash(hash.Hash)
	GetPosition() position
}

func EvaluableSum(in Evaluable) string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	in.WriteHash(hash)
	return fmt.Sprintf("%016x", hash.Sum64())
}

func offsetSources(input []StringSource, offset int) []StringSource {
	result := make([]StringSource, len(input))
	for i, item := range input {
		result[i] = StringSource{item.Start + offset, item.Len, item.Value}
	}
	return result
}

func SourceAt(input []StringSource, pos int) []Evaluable {
	inside := make([]Evaluable, 0, len(input))

	for _, item := range input {
		if pos >= item.Start && pos < item.Start+item.Len {
			inside = append(inside, item.Value)
		}
	}

	return inside
}
