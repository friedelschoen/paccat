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

func (this *recipeImport) HasOutput() bool {
	/* a function call must be at the root of the recipe and a `${out}` in the
	   path or parameters wouldn't make sense, we just say we don't have any output */
	return false
}

func (this *recipeImport) Eval(ctx *Context, attr string) (string, error) {
	if attr == "" {
		attr = DefaultAttribute
	}
	filename, err := this.source.Eval(ctx, "")
	if err != nil {
		return "", err
	}

	pathname := path.Join(ctx.workDir, filename)
	recipe, err := ParseFile(pathname)
	if err != nil {
		return "", err
	}

	newContext, err := recipe.(*Recipe).NewContext(path.Dir(pathname), this.arguments)
	if err != nil {
		return "", err
	}

	value, ok := newContext.scope[attr] //(attr, false)
	if !ok {
		return "", UnknownAttributeError{ctx, this.pos, filename, attr}
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
