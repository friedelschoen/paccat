package recipe

import (
	"fmt"
	"hash"
)

type Recipe struct {
	attributes         map[string]Evaluable
	requiredAttributes map[string]Evaluable
}

func (this *Recipe) String() string {
	return fmt.Sprintf("Recipe{ require=%s, attr=%s }", this.requiredAttributes, this.attributes)
}

func (this *Recipe) NewContext(directory string, params map[string]Evaluable) (*Context, error) {
	ctx := &Context{
		currentRecipe: this,
		scope:         this.attributes,
		workDir:       directory,
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
	for key, value := range this.requiredAttributes {
		hash.Write([]byte(key))
		value.WriteHash(hash)
	}
	for key, value := range this.attributes {
		hash.Write([]byte(key))
		if value != nil {
			value.WriteHash(hash)
		}
	}
}
