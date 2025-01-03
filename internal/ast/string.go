package ast

import (
	"fmt"
	"hash"
	"strings"

	"friedelschoen.io/paccat/internal/errors"
)

type StringNode struct {
	Pos     errors.Position
	Content []Node
}

func (this *StringNode) String() string {
	builder := strings.Builder{}
	builder.WriteString("RecipeString{")
	for i, content := range this.Content {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", content))
	}
	return builder.String()
}

func (this *StringNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("string"))
	for _, content := range this.Content {
		content.WriteHash(hash)
	}
}

func (this *StringNode) GetPosition() errors.Position {
	return this.Pos
}
