package parser

import (
	"io"
	"os"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
)

func Parse(filename, content string) (ast.Node, error) {
	file := &errors.ErrorFile{
		Filename: filename,
		Content:  content,
	}
	parser := parseState{}
	parser.File = file
	parser.current = []state{tokens[0].state}

	parser.Next()
	result, err := parser.parseFile()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ParseFile(filename string) (ast.Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return Parse(filename, string(content))
}
