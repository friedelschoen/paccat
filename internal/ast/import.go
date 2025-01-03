package ast

import (
	"hash"
	"path"
)

const DefaultAttribute = "build"

type recipeImport struct {
	pos    Position
	source Evaluable
}

func (this *recipeImport) String() string {
	return "RecipeImport"
}

func (this *recipeImport) Eval(ctx Context) (Value, error) {
	filenameVal, err := this.source.Eval(ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while evaluating import")
	}
	filename, err := CastString(filenameVal, ctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while evaluating import")
	}

	pathname := path.Join(ctx.workdir, filename.Content)
	recipe, err := ParseFile(pathname)
	if err != nil {
		return nil, err
	}

	newctx := Context{
		workdir: path.Dir(pathname),
	}
	value, err := recipe.(Evaluable).Eval(newctx)
	if err != nil {
		return nil, WrapRecipeError(err, this.pos, "while evaluating import")
	}
	return value, nil
}

func (this *recipeImport) WriteHash(hash hash.Hash) {
	hash.Write([]byte("import"))
	this.source.WriteHash(hash)
}

func (this *recipeImport) GetPosition() Position {
	return this.pos
}
