package types

import (
	"fmt"
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
	MaxSimilarityDistance = 10
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

func (this *Context) findSimilar(name string) (string, int) {
	lowest := ""
	lowestDist := math.MaxInt
	for current := range this.scope {
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

func (ctx *Context) Unwrap(currentNode ast.Node) (ast.Node, *Context, error) {
	for {
		switch this := currentNode.(type) {
		case *ast.ImportNode:
			filename, err := ctx.Evaluate(this.Source)
			if err != nil {
				return nil, nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
			}

			pathname := path.Join(ctx.workdir, filename.Content)
			currentNode, err = parser.ParseFile(pathname)
			if err != nil {
				return nil, nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
			}
			ctx = &Context{
				workdir: path.Dir(pathname),
				scope:   map[string]ast.Node{},
			}
		case *ast.CallNode:
			target, ctx, err := ctx.Unwrap(this.Target)
			if err != nil {
				return nil, nil, errors.WrapRecipeError(err, this.GetPosition(), "unable to call "+this.Target.Name())
			}
			lambda, ok := target.(*ast.LambdaNode)
			if !ok {
				return nil, nil, errors.NewRecipeError(this.GetPosition(), "unable to call "+this.Target.Name())
			}

			for key, def := range lambda.Args {
				if val, ok := this.Args[key]; ok {
					ctx.scope[key] = val.Value
				} else if val.Value != nil {
					ctx.scope[key] = def.Value
				} else {
					return nil, nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("lambda called without parameter `%s`", key))
				}
			}
		default:
			return currentNode, ctx, nil
		}
	}
}

func (ctx *Context) Evaluate(currentNode ast.Node) (*StringValue, error) {
	currentNode, ctx, err := ctx.Unwrap(currentNode)
	if err != nil {
		return nil, err
	}
	switch this := currentNode.(type) {
	case *ast.GetterNode:
		value, err := ctx.Evaluate(this.Target)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), fmt.Sprintf("while trying to get attribute `%s`", this.Attribute.Content))
		}

		res, ok := value.Attributes[this.Attribute.Content]
		if !ok {
			return nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("while trying to get attribute `%s`", this.Attribute.Content))
		}
		return res.Value, nil
	case *ast.DictNode:
		values := map[string]ValuePair{}
		for key, itempair := range this.Items {
			pair := ValuePair{}
			pair.Key, err = ctx.Evaluate(itempair.Key)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluting dict")
			}
			pair.Value, err = ctx.Evaluate(itempair.Value)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluting dict")
			}
			values[key] = pair
		}
		return &StringValue{
			Node:       this,
			Attributes: values,
		}, nil
	case *ast.ListNode:
		builder := ValueBuilder{}
		for i, item := range this.Items {
			if i > 0 {
				builder.WriteByte(' ')
			}
			anyValue, err := ctx.Evaluate(item)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.Pos, "while evaluating list")
			}
			builder.WriteValue(anyValue)
		}
		return builder.Value(this), nil
	case *ast.LiteralNode:
		return &StringValue{
			Node:    this,
			Content: this.Content,
		}, nil
	case *ast.OutputNode:
		for key, value := range this.Options {
			ctx.scope[key] = value.Value
		}

		sum := this.ScriptSum()
		outpath := path.Join(util.GetCachedir(), sum)

		if _, err := os.Stat(outpath); err == nil {
			if alwaysEval, ok := ctx.scope["always"]; ok {
				alwaysVal, err := ctx.Evaluate(alwaysEval)
				if err != nil {
					return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
				}
				if len(alwaysVal.Content) > 0 {
					return &StringValue{
						Node:    this,
						Content: outpath,
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

		ctx.Set("out", outpath)
		defer ctx.Unset("out")

		if scriptEval, ok := ctx.scope["script"]; ok {
			scriptValue, err := ctx.Evaluate(scriptEval)
			if err != nil {
				return nil, errors.WrapRecipeError(err, scriptEval.GetPosition(), "while evaluating output")
			}

			cmd := exec.Command("sh")
			cmd.Stdin = strings.NewReader(scriptValue.Content)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = workdir
			if err = cmd.Run(); err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
			}
		}

		return &StringValue{
			Node:    this,
			Content: outpath,
		}, nil
	case *ast.PanicNode:
		value, err := ctx.Evaluate(this.Message)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating panic")
		}

		return nil, errors.NewRecipeError(this.GetPosition(), value.Content)
	case *ast.ReferenceNode:
		value, ok := ctx.scope[this.Variable.Content]
		if !ok {
			if len(ctx.scope) > 0 {
				similar, dist := ctx.findSimilar(this.Variable.Content)
				if dist <= MaxSimilarityDistance {
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
		builder := ValueBuilder{}
		for _, content := range this.Content {
			value, err := ctx.Evaluate(content)
			if err != nil {
				return nil, err
			}
			builder.WriteValue(value)
		}
		return builder.Value(this), nil
	default:
		return nil, errors.NewRecipeError(currentNode.GetPosition(), currentNode.Name()+" is not evaluable")
	}
}
