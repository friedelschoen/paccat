package recipe

import "fmt"

type RecipeGetter struct {
	reference Buildable
	attribute string
}

func (this *RecipeGetter) String() string {
	return fmt.Sprintf("RecipeGetter#%s{%v}", this.attribute, this.reference)
}

func (this *RecipeGetter) HasOutput() bool {
	return this.reference.HasOutput()
}

func (this *RecipeGetter) Build(ctx *Context) (string, error) {
	ctx.nextAttribute = &this.attribute
	value, err := this.reference.Build(ctx)
	if err != nil {
		return "", err
	}
	if ctx.nextAttribute != nil {
		return "", fmt.Errorf("attribute-getter not applied on function")
	}
	return value, nil
}
