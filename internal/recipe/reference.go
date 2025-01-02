package recipe

import (
	"fmt"
	"hash"
	"iter"
	"maps"
	"math"

	"github.com/agnivade/levenshtein"
)

const (
	SimilarDistance = 10
)

type recipeReference struct {
	pos  Position
	name string
}

func findSimilar(name string, vars iter.Seq[string]) (string, int) {
	lowest := ""
	lowestDist := math.MaxInt
	for current := range vars {
		if dist := levenshtein.ComputeDistance(name, current); dist < lowestDist {
			lowest = current
			lowestDist = dist
		}
	}
	return lowest, lowestDist
}

func (this *recipeReference) Eval(ctx Context) (Value, error) {
	value, ok := ctx.scope[this.name]
	if !ok {
		if len(ctx.scope) > 0 {
			similar, dist := findSimilar(this.name, maps.Keys(ctx.scope))
			if dist <= SimilarDistance {
				return nil, NewRecipeError(this.GetPosition(), fmt.Sprintf("`%s` is not defined in current scope, do you mean `%s`?", this.name, similar))
			}
		}
		return nil, NewRecipeError(this.GetPosition(), fmt.Sprintf("`%s` is not defined in current scope", this.name))
	}
	eval, err := value.Eval(ctx)
	if err != nil {
		WrapRecipeError(err, this.GetPosition(), "refered here")
	}
	return eval, nil
}

func (this *recipeReference) String() string {
	return fmt.Sprintf("RecipeReference#%s", this.name)
}

func (this *recipeReference) WriteHash(hash hash.Hash) {
	hash.Write([]byte("reference"))
	hash.Write([]byte(this.name))
}

func (this *recipeReference) GetPosition() Position {
	return this.pos
}
