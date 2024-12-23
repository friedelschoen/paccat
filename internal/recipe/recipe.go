package recipe

import (
	"fmt"
	"hash"
	"path"
)

type Recipe struct {
	pos                position
	attributes         map[string]Evaluable
	requiredAttributes map[string]Evaluable
}

func (this *Recipe) String() string {
	return fmt.Sprintf("Recipe{ require=%s, attr=%s }", this.requiredAttributes, this.attributes)
}

func (this *Recipe) NewContext(filename string, params map[string]Evaluable) (*Context, error) {
	ctx := &Context{
		currentRecipe: this,
		scope:         this.attributes,
		workDir:       path.Dir(filename),
		filename:      path.Base(filename),
		forceBuild:    false,
	}

	/* override attributes */
	if params != nil {
		for key, value := range params {
			ctx.scope[key] = value
		}
	}

	for key, value := range this.requiredAttributes {
		_, ok := ctx.scope[key]
		if !ok {
			if value == nil {
				return nil, fmt.Errorf("recipe requires key: %s", key)
			}
			ctx.scope[key] = value
		}
	}

	return ctx, nil
}

func (this *Recipe) WriteHash(hash hash.Hash) {
	hash.Write([]byte("recipe"))
	writeHashMap(this.attributes, hash)
	writeHashMap(this.requiredAttributes, hash)
}

func (this *Recipe) GetPosition() position {
	return this.pos
}
