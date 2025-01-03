package types

import (
	"friedelschoen.io/paccat/internal/ast"
)

type OutputValue struct {
	Source  ast.Node
	Path    string              /* outdir */
	Exports map[string][]string /* environment variables, for example PATH=.../bin */
}

func (this *OutputValue) GetSource() ast.Node {
	return this.Source
}

func (this *OutputValue) GetName() string {
	return "output"
}

func (this *OutputValue) ToString(ctx Context) (*StringValue, error) {
	return &StringValue{
		source:       this.Source,
		Content:      this.Path,
		StringSource: []StringSource{},
	}, nil
}
