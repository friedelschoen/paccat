package recipe

import (
	"fmt"
	"hash"
	"path"
)

type recipeImport struct {
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

func (this *recipeImport) Eval(ctx *Context) (string, error) {
	if ctx.importAttribute == nil && !ctx.isBuilding {
		return "", NoGetterError
	}

	filename, err := this.source.Eval(ctx)
	if err != nil {
		return "", err
	}

	pathname := path.Join(ctx.workDir, filename)
	recipe, err := ParseFile(pathname)
	if err != nil {
		return "", err
	}

	newContex, err := recipe.(*Recipe).NewContext(path.Dir(pathname), this.arguments)
	if err != nil {
		return "", err
	}

	if ctx.importAttribute == nil {
		return newContex.EvalPackage()
	} else {
		attr := *ctx.importAttribute
		ctx.importAttribute = nil

		return newContex.Get(attr, false)
	}
}

func (this *recipeImport) WriteHash(hash hash.Hash) {
	this.source.WriteHash(hash)
	for key, value := range this.arguments {
		hash.Write([]byte(key))
		value.WriteHash(hash)
	}
}
