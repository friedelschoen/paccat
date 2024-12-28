package recipe

import (
	"fmt"
	"hash"
	"path"
)

const DefaultAttribute = "build"

type recipeImport struct {
	pos       position
	source    Evaluable
	arguments map[string]Evaluable
}

func (this *recipeImport) String() string {
	return fmt.Sprintf("RecipeImport#%v{%v}", this.source, this.arguments)
}

func (this *recipeImport) Eval(ctx *Context, attr string) (string, []StringSource, error) {
	if attr == "" {
		attr = DefaultAttribute
	}
	filename, _, err := this.source.Eval(ctx, "")
	if err != nil {
		return "", nil, err
	}

	pathname := path.Join(ctx.workDir, filename)
	recipe, err := ParseFile(pathname)
	if err != nil {
		return "", nil, err
	}

	newContext, err := recipe.(*Recipe).NewContext(path.Dir(pathname), this.arguments, ctx.Database)
	if err != nil {
		return "", nil, err
	}

	value, ok := newContext.scope[attr] //(attr, false)
	if !ok {
		return "", nil, UnknownAttributeError{ctx, this.pos, filename, attr}
	}
	return value.Eval(newContext, "")
}

func (this *recipeImport) WriteHash(hash hash.Hash) {
	hash.Write([]byte("import"))
	this.source.WriteHash(hash)
	writeHashMap(this.arguments, hash)
}

func (this *recipeImport) GetPosition() position {
	return this.pos
}
