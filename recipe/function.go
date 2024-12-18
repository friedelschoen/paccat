package recipe

import (
	"errors"
	"fmt"
	"os"
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
	if ctx.nextAttribute == nil {
		return "", errors.New("function mentioned without getter")
	}

	attr := *ctx.nextAttribute
	fmt.Printf("attr: %s!\n", attr)
	ctx.nextAttribute = nil

	filename, err := this.Path.Build(ctx)
	if err != nil {
		return "", err
	}

	path := path.Join(ctx.directory, filename)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	got, err := ParseRecipe(path, file)
	if err != nil {
		return "", err
	}

	return got.Build(attr, true, this.Parameters)
}
