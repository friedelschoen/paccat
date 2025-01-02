package recipe

import (
	"fmt"
	"hash"
	"hash/crc64"
	"os"
	"os/exec"
	"path"
	"strings"

	"friedelschoen.io/paccat/internal/util"
)

type recipeOutput struct {
	pos     Position
	options map[string]Evaluable
}

type outputOption func(*recipeOutput)

func (this *recipeOutput) String() string {
	return fmt.Sprintf("RecipeOutput{%v}", this.options)
}

func (this *recipeOutput) WriteHash(hash hash.Hash) {
	hash.Write([]byte("output"))
	for key, value := range this.options {
		hash.Write([]byte(key))
		value.WriteHash(hash)
	}
}

func appendEnv(env []string, key string, value string) []string {
	key += "="
	for i, pair := range env {
		if strings.HasPrefix(pair, key) {
			env[i] = pair + ":" + value
			return env
		}
	}
	return append(env, key+value)
}

func (this *recipeOutput) MakeEnviron(ctx *Context) []string {
	return os.Environ()
}

func (this *recipeOutput) ScriptSum() string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	for key, value := range this.options {
		hash.Write([]byte(key))
		value.WriteHash(hash)
	}
	return fmt.Sprintf("%016x", hash.Sum64())
}

func (this *recipeOutput) Exports(ctx Context) (map[string][]string, error) {
	exports := map[string][]string{}
	if exportVal, ok := ctx.scope["exports"]; ok {
		anyValue, err := exportVal.Eval(ctx)
		if err != nil {
			return nil, WrapRecipeError(err, this.pos, "while evaluating output")
		}
		switch value := anyValue.(type) {
		case *recipeDict:
			for key, attr := range value.items {
				attrValue, err := attr.Eval(ctx)
				if err != nil {
					return nil, WrapRecipeError(err, this.pos, "while evaluating output")
				}
				if listval, ok := attrValue.(*recipeList); ok {
					vars := make([]string, len(listval.items))
					for i, item := range listval.items {
						itemValue, err := item.Eval(ctx)
						if err != nil {
							return nil, WrapRecipeError(err, this.pos, "while evaluating output")
						}
						strval, err := CastString(itemValue, ctx)
						if err != nil {
							return nil, WrapRecipeError(err, this.pos, "while evaluating output")
						}
						vars[i] = strval.Content
					}
					exports[key] = vars
				} else {
					strval, err := CastString(attrValue, ctx)
					if err != nil {
						return nil, WrapRecipeError(err, this.pos, "while evaluating output")
					}
					exports[key] = strings.Split(strval.Content, ":")
				}
			}

		default:
			return nil, NewRecipeError(this.pos, "option `exports` is not a list or dict")
		}
	}
	return exports, nil
}

func (this *recipeOutput) Eval(oldctx Context) (Value, error) {
	ctx := oldctx.Copy()
	for key, value := range this.options {
		ctx.scope[key] = value
	}

	sum := this.ScriptSum()
	outpath := path.Join(util.GetCachedir(), sum)

	exports, err := this.Exports(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(outpath); err == nil {
		if alwaysEval, ok := ctx.scope["always"]; ok {
			alwaysVal, err := alwaysEval.Eval(ctx)
			if err != nil {
				return nil, WrapRecipeError(err, this.pos, "while evaluating output")
			}
			always, err := CastBoolean(alwaysVal, ctx)
			if err != nil {
				return nil, WrapRecipeError(err, alwaysEval.GetPosition(), "while evaluating output")
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
			return nil, WrapRecipeError(err, this.pos, "while cleaning output")
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
		scriptValue, err := scriptEval.Eval(ctx)
		if err != nil {
			return nil, WrapRecipeError(err, scriptEval.GetPosition(), "while evaluating output")
		}
		script, err := CastString(scriptValue, ctx)
		if err != nil {
			return nil, WrapRecipeError(err, scriptEval.GetPosition(), "while evaluating output")
		}

		cmd := exec.Command("sh")
		cmd.Stdin = strings.NewReader(script.Content)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = workdir
		if err = cmd.Run(); err != nil {
			return nil, WrapRecipeError(err, this.pos, "while evaluating output")
		}
	}

	return &OutputValue{
		Source:  this,
		Path:    outpath,
		Exports: exports,
	}, nil
}

func (this *recipeOutput) GetPosition() Position {
	return this.pos
}

type OutputValue struct {
	Source  Evaluable
	Path    string              /* outdir */
	Exports map[string][]string /* environment variables, for example PATH=.../bin */
}

func (this *OutputValue) GetSource() Evaluable {
	return this.Source
}

func (this *OutputValue) GetName() string {
	return "output"
}

func (this *OutputValue) ToString(ctx Context) (*StringValue, error) {
	return &StringValue{
		source:       this.Source,
		Content:      this.Path,
		StringSource: []StringSource{},
	}, nil
}
