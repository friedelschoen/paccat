package ast

import (
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"unicode"

	"friedelschoen.io/paccat/internal/errors"
)

type Positioned interface {
	GetPosition() errors.Position
}

type Node interface {
	Positioned
	Name() string
	GetChildren() []Node
	WriteHash(hash.Hash)
}

func NodeHash(in Node) string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	in.WriteHash(hash)
	return fmt.Sprintf("%016x", hash.Sum64())
}

func PrintTree(w io.Writer, node Node, level int) {
	indent := []byte("    ")

	pos := node.GetPosition()
	name := []rune(node.Name())
	for i, chr := range name {
		if !unicode.IsGraphic(chr) {
			name[i] = '-'
		}
	}
	fmt.Fprintf(w, "%s at %d-%d: %s\n", string(name), pos.Start, pos.End, NodeHash(node))

	for _, child := range node.GetChildren() {
		for ; level > 0; level-- {
			w.Write(indent)
		}
		w.Write([]byte("- "))
		PrintTree(w, child, level+1)
	}
}
