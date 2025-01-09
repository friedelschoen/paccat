package ast

import (
	"math"

	"friedelschoen.io/paccat/internal/errors"
)

type LiteralNode struct {
	Pos     errors.Position
	Content string
}

func (this *LiteralNode) Name() string {
	return "literal-" + this.Content
}

func (this *LiteralNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *LiteralNode) GetChildren() []Node {
	return []Node{}
}

type LiteralMap map[string]LiteralMapPair

type LiteralMapPair struct {
	Key   *LiteralNode
	Value Node
}

func (this LiteralMap) Name() string {
	return "literalmap"
}

func (this LiteralMap) GetPosition() errors.Position {
	pos := errors.Position{}
	pos.Start = math.MaxInt
	pos.End = 0
	for _, pair := range this {
		if pos.Content == nil {
			pos.Filename = pair.Key.Pos.Filename
			pos.Content = pair.Key.Pos.Content
		}
		if pair.Key.Pos.Start < pos.Start {
			pos.Start = pair.Key.Pos.Start
		}
		if pair.Key.Pos.End > pos.End {
			pos.End = pair.Key.Pos.End
		}
		if pair.Value.GetPosition().Start < pos.Start {
			pos.Start = pair.Value.GetPosition().Start
		}
		if pair.Value.GetPosition().End > pos.End {
			pos.End = pair.Value.GetPosition().End
		}
	}
	return pos
}

func (this LiteralMap) GetChildren() []Node {
	res := make([]Node, 2*len(this))
	i := 0
	for _, pair := range this {
		res[i] = pair.Key
		i++
		res[i] = pair.Value
		i++
	}
	return res
}
