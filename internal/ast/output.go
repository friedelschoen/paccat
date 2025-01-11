package ast

import (
	"strings"

	"friedelschoen.io/paccat/internal/errors"
)

type OutputNode struct {
	Pos     errors.Position
	Options Node
}

func (this *OutputNode) Name() string {
	return "output"
}

func appendEnv(env []string, key string, value string) []string {
	key += "="
	for i, pair := range env {
		if strings.HasPrefix(pair, key) {
			env[i] = pair + ":" + value
			return env
		}
	}
	return append(env, key+value)
}

func (this *OutputNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *OutputNode) GetChildren() []Node {
	return []Node{this.Options}
}
