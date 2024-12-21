package recipe

import "hash"

type Evaluable interface {
	Eval(*Context) (string, error)
	HasOutput() bool
	WriteHash(hash.Hash)
}
