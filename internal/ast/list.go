package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type ListNode struct {
	Pos   errors.Position
	Items []Node
}

func (this *ListNode) String() string {
	return "RecipeList"
}

func (this *ListNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	for _, value := range this.Items {
		value.WriteHash(hash)
	}
}

func (this *ListNode) GetPosition() errors.Position {
	return this.Pos
}
