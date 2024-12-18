package recipe

import (
	"fmt"
)

type Recipe struct {
	Attributes, RequireAttributes map[string]Buildable
}

func (this *Recipe) String() string {
	return fmt.Sprintf("Recipe{ require=%s, attr=%s }", this.RequireAttributes, this.Attributes)
}

func (this *Recipe) Build(buildAttr string, forceOutput bool, params map[string]Buildable) (string, error) {
	ctx := Context{
		currentRecipe: this,
		attributes:    this.Attributes,
	}

	/* override attributes */
	for key, value := range params {
		ctx.attributes[key] = value
	}

	for key, value := range this.RequireAttributes {
		_, ok := ctx.attributes[key]
		if !ok {
			if value == nil {
				return "", fmt.Errorf("recipe requires key: %s", key)
			}
			ctx.attributes[key] = value
		}
	}

	return ctx.Get(buildAttr, true)
}
