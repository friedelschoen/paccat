package recipe

import (
	"errors"
	"fmt"
	"hash"
	"os"
	"os/exec"
	"path"
	"strings"

	"friedelschoen.io/paccat/internal/util"
)

type recipeOutput struct {
	pos    position
	script Evaluable
}

func (this *recipeOutput) String() string {
	return fmt.Sprintf("RecipeOutput{%v}", this.script)
}

func (this *recipeOutput) HasOutput() bool {
	return true
}

func (this *recipeOutput) WriteHash(hash hash.Hash) {
	hash.Write([]byte("output"))
	this.script.WriteHash(hash)
}

func createOutDir(hash uint64) string {
	name := fmt.Sprintf("%16x", hash)
	return path.Join(util.GetCachedir(), name)
}

func (this *recipeOutput) Eval(ctx *Context, attr string) (string, error) {
	if attr != "" {
		return "", NoAttributeError{ctx, this.pos, "output-statement", attr}
	}

	sum := EvaluableSum(this.script)
	outdir := createOutDir(sum)

	if _, err := os.Stat(outdir); errors.Is(err, os.ErrExist) {
		if !ctx.forceBuild {
			return outdir, nil
		}
		if err = os.RemoveAll(outdir); err != nil {
			return "", err
		}
	}

	workdir, err := os.MkdirTemp(os.TempDir(), "paccat-workdir-******")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workdir) /* do remove the workdir if not needed */

	ctx.scope["out"] = &recipeStringLiteral{position{}, outdir}
	defer delete(ctx.scope, "out")

	script, err := this.script.Eval(ctx, "")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh")
	cmd.Stdin = strings.NewReader(script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workdir
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return outdir, nil
}

func (this *recipeOutput) GetPosition() position {
	return this.pos
}
