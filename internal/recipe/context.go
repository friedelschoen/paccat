package recipe

type Context struct {
	filename      string
	workDir       string               // directory of the recipe
	currentRecipe *Recipe              // current recipe
	scope         map[string]Evaluable // variables and attributes
	forceBuild    bool                 // build output if directory exists
}

func (this *Context) AlwaysBuild() {
	this.forceBuild = true
}

func (this *Context) Set(key, value string) {
	this.scope[key] = &recipeStringLiteral{position{}, value}
}

func (this *Context) Unset(key string) {
	delete(this.scope, key)
}

func (this *Context) Get(name, attr string) (string, error) {
	value, ok := this.scope[name]
	if !ok {
		return "", UnknownReferenceError{this, position{}, name}
	}
	return value.Eval(this, attr)
}
