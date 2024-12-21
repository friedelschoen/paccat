package recipe

import (
	"fmt"
	"hash/maphash"
	"path"
)

type RecipeFunction struct {
	Path       Buildable
	Parameters map[string]Buildable
}

func (this *RecipeFunction) String() string {
	return fmt.Sprintf("RecipeFunction#%v{%v}", this.Path, this.Parameters)
}

func (this *RecipeFunction) HasOutput() bool {
	/* a function call must be at the root of the recipe and a `${out}` in the
	   path or parameters wouldn't make sense, we just say we don't have any output */
	return false
}

func (this *RecipeFunction) Build(ctx *Context) (string, error) {
	if ctx.nextAttribute == nil && !ctx.building {
		return "", NoGetterError
	}

	filename, err := this.Path.Build(ctx)
	if err != nil {
		return "", err
	}

	pathname := path.Join(ctx.directory, filename)
	recipe, err := ParseFile(pathname)
	if err != nil {
		return "", err
	}

	newContex, err := recipe.(*Recipe).NewContext(path.Dir(pathname), this.Parameters)
	if err != nil {
		return "", err
	}

	if ctx.nextAttribute == nil {
		return newContex.BuildPackage()
	} else {
		attr := *ctx.nextAttribute
		ctx.nextAttribute = nil

		return newContex.Get(attr, false)
	}
}

func (this *RecipeFunction) WriteHash(hash maphash.Hash) {
	this.Path.WriteHash(hash)
	for key, value := range this.Parameters {
		hash.WriteString(key)
		value.WriteHash(hash)
	}
}
