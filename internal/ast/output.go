package ast

import (
	"fmt"
	"hash"
	"hash/crc64"
	"strings"

	"friedelschoen.io/paccat/internal/errors"
)

type OutputNode struct {
	Pos     errors.Position
	Options LiteralMap
}

func (this *OutputNode) Name() string {
	return "output"
}

func (this *OutputNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("output"))
	this.Options.WriteHash(hash)
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

func (this *OutputNode) ScriptSum() string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	this.WriteHash(hash)
	return fmt.Sprintf("%016x", hash.Sum64())
}

func (this *OutputNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *OutputNode) GetChildren() []Node {
	return []Node{this.Options}
}
