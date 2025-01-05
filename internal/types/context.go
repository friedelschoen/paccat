package types

import (
	"fmt"
	"iter"
	"maps"
	"math"
	"os"
	"os/exec"
	"path"
	"strings"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
	"friedelschoen.io/paccat/internal/parser"
	"friedelschoen.io/paccat/internal/util"
	"github.com/agnivade/levenshtein"
)

const (
	SimilarDistance = 10
)

type Context struct {
	workdir string              // directory of the recipe
	scope   map[string]ast.Node // variables and attributes
}

func NewContext(filename string) Context {
	return Context{
		workdir: path.Dir(filename),
		scope:   map[string]ast.Node{},
	}
}

func GetExports(ctx Context, this *ast.OutputNode) (map[string][]string, error) {
	exports := map[string][]string{}
	if exportVal, ok := ctx.scope["exports"]; ok {
		anyValue, err := ctx.Evaluate(exportVal)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
		}
		switch value := anyValue.(type) {
		case *DictValue:
			for key, pair := range value.Items {
				attrValue, err := ctx.Evaluate(pair.Value)
				if err != nil {
					return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
				}
				if listval, ok := attrValue.(*ListValue); ok {
					vars := make([]string, len(listval.Items))
					for i, item := range listval.Items {
						itemValue, err := ctx.Evaluate(item)
						if err != nil {
							return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
						}
						strval, err := CastString(itemValue, ctx)
						if err != nil {
							return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
						}
						vars[i] = strval.Content
					}
					exports[key] = vars
				} else {
					strval, err := CastString(attrValue, ctx)
					if err != nil {
						return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
					}
					exports[key] = strings.Split(strval.Content, ":")
				}
			}

		default:
			return nil, errors.NewRecipeError(this.GetPosition(), "option `exports` is not a list or dict")
		}
	}
	return exports, nil
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

func (this *Context) Copy() Context {
	newctx := Context{
		workdir: this.workdir,
		scope:   map[string]ast.Node{},
	}

	for key, value := range this.scope {
		newctx.scope[key] = value
	}

	return newctx
}

func (this *Context) Set(key, value string) {
	literal := fmt.Sprintf("\"%s\"", value)
	this.scope[key] = &ast.LiteralNode{
		Pos:     errors.Position{Filename: "<eval>", Content: &literal, Start: 0, End: len(literal)},
		Content: value,
	}
}

func (this *Context) Unset(key string) {
	delete(this.scope, key)
}

func (this *Context) Hash(name string) bool {
	_, ok := this.scope[name]
	return ok
}

func (ctx Context) Evaluate(anyNode ast.Node) (Value, error) {
	switch this := anyNode.(type) {
	case *ast.GetterNode:
		anyValue, err := ctx.Evaluate(this.Target)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), fmt.Sprintf("while trying to get attribute `%s`", this.Attribute.Content))
		}

		dict, ok := anyValue.(DictLike)
		if !ok {
			return nil, errors.NewRecipeError(anyValue.GetSource().GetPosition(), fmt.Sprintf("cannot cast %s to dict", anyValue.GetName()))
		}

		res, err := dict.GetAttrbute(this.Attribute.Content, ctx)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), fmt.Sprintf("while trying to get attribute `%s`", this.Attribute.Content))
		}
		return res, nil
	case *ast.CallNode:
		value, err := ctx.Evaluate(this.Target)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while trying to call value")
		}
		lambda, err := CastValue[*LambdaValue](value)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while trying to call value")
		}

		newctx := ctx.Copy()
		for key, def := range lambda.Args {
			if val, ok := this.Args[key]; ok {
				newctx.scope[key] = val.Value
			} else if val.Value != nil {
				newctx.scope[key] = def.Value
			} else {
				return nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("lambda called without parameter `%s`", key))
			}
		}

		res, err := newctx.Evaluate(lambda.Target)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while trying to call value")
		}
		return res, nil
	case *ast.DictNode:
		return (*DictValue)(this), nil
	case *ast.ImportNode:
		filenameVal, err := ctx.Evaluate(this.Source)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
		}
		filename, err := CastString(filenameVal, ctx)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
		}

		pathname := path.Join(ctx.workdir, filename.Content)
		recipe, err := parser.ParseFile(pathname)
		if err != nil {
			return nil, err
		}

		newctx := Context{
			workdir: path.Dir(pathname),
		}
		value, err := newctx.Evaluate(recipe)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
		}
		return value, nil
	case *ast.LambdaNode:
		return (*LambdaValue)(this), nil
	case *ast.ListNode:
		return (*ListValue)(this), nil
	case *ast.LiteralNode:
		return &StringValue{
			source:       this,
			Content:      this.Content,
			StringSource: []StringSource{},
		}, nil
	case *ast.OutputNode:
		newctx := ctx.Copy()
		for key, value := range this.Options {
			newctx.scope[key] = value.Value
		}

		sum := this.ScriptSum()
		outpath := path.Join(util.GetCachedir(), sum)

		exports, err := GetExports(newctx, this)
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat(outpath); err == nil {
			if alwaysEval, ok := newctx.scope["always"]; ok {
				alwaysVal, err := newctx.Evaluate(alwaysEval)
				if err != nil {
					return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
				}
				always, err := CastBoolean(alwaysVal, newctx)
				if err != nil {
					return nil, errors.WrapRecipeError(err, alwaysEval.GetPosition(), "while evaluating output")
				}
				if always {
					return &OutputValue{
						Source:  this,
						Exports: exports,
						Path:    outpath,
					}, nil
				}
			}
			if err = os.RemoveAll(outpath); err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while cleaning output")
			}
		}

		workdir, err := os.MkdirTemp(os.TempDir(), "paccat-workdir-")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(workdir) /* do remove the workdir if not needed */

		newctx.Set("out", outpath)
		defer newctx.Unset("out")

		if scriptEval, ok := newctx.scope["script"]; ok {
			scriptValue, err := newctx.Evaluate(scriptEval)
			if err != nil {
				return nil, errors.WrapRecipeError(err, scriptEval.GetPosition(), "while evaluating output")
			}
			script, err := CastString(scriptValue, newctx)
			if err != nil {
				return nil, errors.WrapRecipeError(err, scriptEval.GetPosition(), "while evaluating output")
			}

			cmd := exec.Command("sh")
			cmd.Stdin = strings.NewReader(script.Content)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = workdir
			if err = cmd.Run(); err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
			}
		}

		return &OutputValue{
			Source:  this,
			Path:    outpath,
			Exports: exports,
		}, nil
	case *ast.PanicNode:
		value, err := ctx.Evaluate(this.Message)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating panic")
		}
		strValue, err := CastString(value, ctx)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating panic")
		}

		return nil, errors.NewRecipeError(this.GetPosition(), strValue.Content)
	case *ast.ReferenceNode:
		value, ok := ctx.scope[this.Variable.Content]
		if !ok {
			if len(ctx.scope) > 0 {
				similar, dist := findSimilar(this.Variable.Content, maps.Keys(ctx.scope))
				if dist <= SimilarDistance {
					return nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("`%s` is not defined in current scope, do you mean `%s`?", this.Variable.Content, similar))
				}
			}
			return nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("`%s` is not defined in current scope", this.Variable.Content))
		}
		eval, err := ctx.Evaluate(value)
		if err != nil {
			errors.WrapRecipeError(err, this.GetPosition(), "refered here")
		}
		return eval, nil
	case *ast.StringNode:
		builder := strings.Builder{}
		sources := []StringSource{}
		for _, content := range this.Content {
			value, err := ctx.Evaluate(content)
			if err != nil {
				return nil, err
			}
			strValue, err := CastString(value, ctx)
			if err != nil {
				return nil, err
			}
			sources = append(sources, StringSource{builder.Len(), len(strValue.Content), strValue})
			builder.WriteString(strValue.Content)
		}
		return &StringValue{this, builder.String(), sources}, nil
	}
	return nil, nil
}
