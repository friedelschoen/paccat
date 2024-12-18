package recipe

import (
	"io"
)

//go:generate pigeon -o parser.go recipe.peg

func wrapPairs(pairs []pair) (map[string]Buildable, error) {
	result := make(map[string]Buildable)

	for _, keyvalue := range pairs {
		if keyvalue.value == nil {
			result[string(keyvalue.key)] = nil
		} else {
			val, err := keyvalue.value.Wrap()
			if err != nil {
				return nil, err
			}
			result[string(keyvalue.key)] = val
		}
	}
	return result, nil
}

func (this key) Wrap() (Buildable, error) {
	return &RecipeReference{string(this)}, nil
}

func (this recipeLiteral) Wrap() (Buildable, error) {
	return &RecipeStringLiteral{string(this)}, nil
}

func (value recipeString) Wrap() (Buildable, error) {
	result := &RecipeString{}
	result.Elements = make([]Buildable, 0)
	for _, element := range value {
		wrapped, err := element.(Wrappable).Wrap()
		if err != nil {
			return nil, nil
		}
		result.Elements = append(result.Elements, wrapped)
	}
	return result, nil
}

func (value functionCall) Wrap() (Buildable, error) {
	result := &RecipeFunction{}
	wrapped, err := value.path.Wrap()
	if err != nil {
		return nil, err
	}
	result.Path = wrapped

	params, err := wrapPairs(value.params)
	if err != nil {
		return nil, err
	}

	result.Parameters = params
	return result, nil
}

func (value attrGetter) Wrap() (Buildable, error) {
	ref, err := value.call.Wrap()
	if err != nil {
		return nil, err
	}
	return &RecipeGetter{
		attribute: string(value.attr),
		reference: ref,
	}, nil
}

func ParseRecipe(name string, reader io.Reader) (*Recipe, error) {
	result, err := ParseReader(name, reader)
	if err != nil {
		return nil, err
	}
	ast := result.(recipe)

	reqAttrs, err := wrapPairs(ast.require)
	if err != nil {
		return nil, err
	}

	attrs, err := wrapPairs(ast.values)
	if err != nil {
		return nil, err
	}

	return &Recipe{
		RequireAttributes: reqAttrs,
		Attributes:        attrs,
	}, nil
}
