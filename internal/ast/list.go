package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type ListNode struct {
	Pos   errors.Position
	Items []Node
}

func (this *ListNode) Name() string {
	return "list"
}

func (this *ListNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *ListNode) GetChildren() []Node {
	return this.Items
}
