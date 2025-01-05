package types

import (
	"fmt"
	"io"

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

func (this *OutputValue) ToJSON(ctx Context, w io.Writer) error {
	fmt.Fprintf(w, "{\"path\":\"%s\",\"exports\":{", this.Path)

	i := 0
	for key, node := range this.Exports {
		if i > 0 {
			w.Write([]byte{','})
		}
		fmt.Fprintf(w, "\"%s\":[", key)
		for j, val := range node {
			if j > 0 {
				w.Write([]byte{','})
			}
			fmt.Fprintf(w, "\"%s\"", val)
		}
		w.Write([]byte{']'})
		i++
	}
	w.Write([]byte("}}"))

	return nil
}
