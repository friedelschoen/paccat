package recipe

import (
	"fmt"
	"path"
)

type Context struct {
	workdir string               // directory of the recipe
	scope   map[string]Evaluable // variables and attributes
}

func NewContext(filename string) Context {
	return Context{
		workdir: path.Dir(filename),
		scope:   map[string]Evaluable{},
	}
}

func (this *Context) Copy() Context {
	newctx := Context{
		workdir: this.workdir,
		scope:   map[string]Evaluable{},
	}

	for key, value := range this.scope {
		newctx.scope[key] = value
	}

	return newctx
}

func (this *Context) Set(key, value string) {
	literal := fmt.Sprintf("\"%s\"", value)
	this.scope[key] = &recipeStringLiteral{Position{"<eval>", &literal, 0, len(literal)}, value}
}

func (this *Context) Unset(key string) {
	delete(this.scope, key)
}

func (this *Context) Hash(name string) bool {
	_, ok := this.scope[name]
	return ok
}
