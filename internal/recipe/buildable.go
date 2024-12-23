package recipe

import (
	"hash"
	"hash/crc64"
)

type Evaluable interface {
	Eval(*Context, string) (string, error)
	HasOutput() bool
	WriteHash(hash.Hash)
	GetPosition() position
}

func EvaluableSum(in Evaluable) uint64 {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	in.WriteHash(hash)
	return hash.Sum64()
}
