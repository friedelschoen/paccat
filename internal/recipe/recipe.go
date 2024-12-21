package recipe

import (
	"fmt"
	"hash/maphash"
)

type Recipe struct {
	Attributes, RequireAttributes map[string]Buildable
}

func (this *Recipe) String() string {
	return fmt.Sprintf("Recipe{ require=%s, attr=%s }", this.RequireAttributes, this.Attributes)
}

func (this *Recipe) NewContext(directory string, params map[string]Buildable) (*Context, error) {
	ctx := &Context{
		currentRecipe: this,
		attributes:    this.Attributes,
		directory:     directory,
	}

	/* override attributes */
	if params != nil {
		for key, value := range params {
			ctx.attributes[key] = value
		}
	}

	for key, value := range this.RequireAttributes {
		_, ok := ctx.attributes[key]
		if !ok {
			if value == nil {
				return nil, fmt.Errorf("recipe requires key: %s", key)
			}
			ctx.attributes[key] = value
		}
	}

	return ctx, nil
}

func (this *Recipe) WriteHash(hash maphash.Hash) {
	for key, value := range this.RequireAttributes {
		hash.WriteString(key)
		value.WriteHash(hash)
	}
	for key, value := range this.Attributes {
		hash.WriteString(key)
		if value != nil {
			value.WriteHash(hash)
		}
	}
}
