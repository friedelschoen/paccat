package ast

import (
	"fmt"
	"hash/crc64"
	"io"
	"unicode"

	"friedelschoen.io/paccat/internal/errors"
)

type Node interface {
	GetPosition() errors.Position
	Name() string
	GetChildren() []Node
}

func writeHash(in Node, w io.Writer) {
	for _, child := range in.GetChildren() {
		w.Write([]byte(child.Name()))
		writeHash(child, w)
	}
}

func NodeHash(in Node) string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	writeHash(in, hash)
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
