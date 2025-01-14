package types

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
	"friedelschoen.io/paccat/internal/parser"
	"friedelschoen.io/paccat/internal/util"
	"github.com/agnivade/levenshtein"
)

const (
	MaxSimilarityDistance = 3
)

type Variable struct {
	name string
	node ast.Node
}

type Scope []Variable

func (this Scope) findSimilar(name string) (string, int) {
	lowest := ""
	lowestDist := math.MaxInt
	for _, current := range this {
		if dist := levenshtein.ComputeDistance(name, current.name); dist < lowestDist {
			lowest = current.name
			lowestDist = dist
		}
	}
	return lowest, lowestDist
}

func (ctx Scope) Get(name string) ast.Node {
	for _, variable := range ctx {
		if variable.name == name {
			return variable.node
		}
	}
	return nil
}

func (ctx Scope) Set(name string, value ast.Node) Scope {
	newctx := make([]Variable, 0, len(ctx))
	for _, variable := range ctx {
		if variable.name != name {
			newctx = append(newctx, variable)
		}
	}
	if value != nil {
		newctx = append(newctx, Variable{name, value})
	}
	return newctx
}

func asLiteral(content string) *ast.LiteralNode {
	return &ast.LiteralNode{
		Pos: errors.Position{
			File:  &errors.ErrorFile{Filename: "<eval>", Content: content},
			Start: 0,
			End:   len(content)},
		Content: content,
	}
}

func (ctx Scope) SetLiteral(name string, content string) Scope {
	return ctx.Set(name, asLiteral(content))
}

func (ctx Scope) Unwrap(currentNode ast.Node) (ast.Node, Scope, error) {
	for {
		switch this := currentNode.(type) {
		case *ast.ImportNode:
			filename, err := ctx.Evaluate(this.Source)
			if err != nil {
				return nil, nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
			}

			workdir := path.Dir(this.Pos.File.Filename)
			pathname := path.Join(workdir, filename.Content)
			currentNode, err = parser.ParseFile(pathname)
			if err != nil {
				return nil, nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating import")
			}
			ctx = []Variable{}
		case *ast.CallNode:
			var target ast.Node
			var err error
			target, ctx, err = ctx.Unwrap(this.Target)
			if err != nil {
				return nil, nil, errors.WrapRecipeError(err, this.GetPosition(), "unable to call "+this.Target.Name())
			}
			lambda, ok := target.(*ast.LambdaNode)
			if !ok {
				return nil, nil, errors.NewRecipeError(this.GetPosition(), "unable to call "+this.Target.Name())
			}

			for key, def := range lambda.Args {
				if val, ok := this.Args[key]; ok {
					ctx = ctx.Set(key, val.Value)
				} else if def.Value != nil {
					ctx = ctx.Set(key, def.Value)
				} else {
					return nil, nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("lambda called without parameter `%s`", key))
				}
			}
			currentNode = lambda.Target
		case *ast.ReferenceNode:
			currentNode = ctx.Get(this.Variable.Content)
			if currentNode == nil {
				similar, dist := ctx.findSimilar(this.Variable.Content)
				if dist <= MaxSimilarityDistance {
					return nil, nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("`%s` is not defined in current scope, do you mean `%s`?", this.Variable.Content, similar))
				}
				return nil, nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("`%s` is not defined in current scope", this.Variable.Content))
			}
		default:
			return currentNode, ctx, nil
		}
	}
}

func makeEnviron(deps *StringValue) []string {
	if deps == nil {
		return os.Environ()
	}
	environ := map[string]string{}
	for _, env := range os.Environ() {
		spl := strings.SplitN(env, "=", 2)
		if len(spl) != 2 {
			continue
		}
		environ[spl[0]] = spl[1]
	}
	for content, dep := range deps.Split() {
		if dep == nil {
			continue
		}
		for name, value := range dep.Attributes {
			if prev, ok := environ[name]; ok {
				environ[name] = fmt.Sprintf("%s:%s/%s", prev, content, value.Content)
			} else {
				environ[name] = content + "/" + value.Content
			}
		}
	}
	result := make([]string, len(environ))
	i := 0
	for key, value := range environ {
		result[i] = key + "=" + value
		i++
	}
	return result
}

func (ctx Scope) Evaluate(currentNode ast.Node) (*StringValue, error) {
	currentNode, ctx, err := ctx.Unwrap(currentNode)
	if err != nil {
		return nil, err
	}
	switch this := currentNode.(type) {
	case *ast.GetterNode:
		value, err := ctx.Evaluate(this.Target)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while trying to get attribute")
		}

		attr, err := ctx.Evaluate(this.Attribute)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while trying to get attribute")
		}
		res, ok := value.Attributes[attr.Content]
		if !ok {
			return nil, errors.NewRecipeError(this.GetPosition(), fmt.Sprintf("target has no attribute `%s`", attr.Content))
		}
		return res, nil
	case *ast.DictNode:
		values := map[string]*StringValue{}
		for key, itempair := range this.Items {
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluting dict")
			}
			values[key], err = ctx.Evaluate(itempair.Value)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluting dict")
			}
		}
		return &StringValue{
			Node:       this,
			Attributes: values,
		}, nil
	case *ast.ListNode:
		builder := ValueBuilder{}
		attrs := make(map[string]*StringValue)
		for i, item := range this.Items {
			if i > 0 {
				builder.WriteByte(' ')
			}
			anyValue, err := ctx.Evaluate(item)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.Pos, "while evaluating list")
			}
			builder.WriteValue(anyValue, true)
			istr := strconv.Itoa(i)
			attrs[istr] = anyValue
		}
		res := builder.Value(this)
		res.Attributes = attrs
		return res, nil
	case *ast.LiteralNode:
		return &StringValue{
			Node:    this,
			Content: this.Content,
		}, nil
	case *ast.NumberNode:
		return &StringValue{
			Node:    this,
			Content: this.Content.Content,
		}, nil
	case *ast.OutputNode:
		var scriptEval ast.Node
		if opt, ok := this.Options.(*ast.DictNode); ok {
			for key, value := range opt.Items {
				if key == "script" {
					scriptEval = value.Value
				}
				ctx = ctx.Set(key, value.Value)
			}
			if scriptEval == nil {
				return nil, errors.NewRecipeError(this.GetPosition(), "output requires field `script`")
			}
		} else {
			scriptEval = this.Options
		}

		sum := ast.NodeHash(this)
		outpath := path.Join(util.GetCachedir(), sum)

		if _, err := os.Stat(outpath); err == nil {
			if alwaysEval := ctx.Get("always"); alwaysEval != nil {
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

		ctx = ctx.SetLiteral("out", outpath)

		var deps *StringValue
		if depsNode := ctx.Get("depends"); depsNode != nil {
			deps, err = ctx.Evaluate(depsNode)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating dependencies")
			}
		}

		var exports map[string]*StringValue
		if exportsNode := ctx.Get("exports"); exportsNode != nil {
			exp, err := ctx.Evaluate(exportsNode)
			if err != nil {
				return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating dependencies")
			}
			exports = exp.Attributes
		}

		scriptValue, err := ctx.Evaluate(scriptEval)
		if err != nil {
			return nil, errors.WrapRecipeError(err, scriptEval.GetPosition(), "while evaluating output")
		}

		cmd := exec.Command("sh")
		cmd.Stdin = strings.NewReader(scriptValue.Content)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = makeEnviron(deps)
		cmd.Dir = workdir
		if err = cmd.Run(); err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating output")
		}

		return &StringValue{
			Node:       this,
			Content:    outpath,
			Attributes: exports,
		}, nil
	case *ast.PanicNode:
		value, err := ctx.Evaluate(this.Message)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.GetPosition(), "while evaluating panic")
		}

		return nil, errors.NewRecipeError(this.GetPosition(), value.Content)
	case *ast.StringNode:
		builder := ValueBuilder{}
		for _, content := range this.Content {
			value, err := ctx.Evaluate(content)
			if err != nil {
				return nil, err
			}
			builder.WriteValue(value, false)
		}
		return builder.Value(this), nil
	case *ast.AttrifyNode:
		target, err := ctx.Evaluate(this.Target)
		if err != nil {
			return nil, errors.WrapRecipeError(err, this.Pos, "while attrifying target")
		}
		builder := &ValueBuilder{}
		for key, value := range target.Attributes {
			if builder.Len() > 0 {
				builder.WriteByte(' ')
			}
			builder.WriteString(key)
			builder.WriteByte('=')
			builder.WriteValue(value, true)
		}
		return builder.Value(this), nil
	default:
		return nil, errors.NewRecipeError(currentNode.GetPosition(), currentNode.Name()+" is not evaluable")
	}
}
