package ast

import (
	"fmt"
	"hash"
	"hash/crc64"

	"friedelschoen.io/paccat/internal/errors"
)

type Positioned interface {
	GetPosition() errors.Position
}

type Node interface {
	Positioned
	WriteHash(hash.Hash)
}

func EvaluableSum(in Node) string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	in.WriteHash(hash)
	return fmt.Sprintf("%016x", hash.Sum64())
}
