package recipe

import "hash/maphash"

type Buildable interface {
	Build(*Context) (string, error)
	HasOutput() bool

	WriteHash(maphash.Hash)
}
