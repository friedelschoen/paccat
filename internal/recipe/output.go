package recipe

import (
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

	always      bool
	try         bool
	interpreter Evaluable
}

type outputOption func(*recipeOutput)

func (this *recipeOutput) String() string {
	return fmt.Sprintf("RecipeOutput{%v}", this.script)
}

func (this *recipeOutput) HasOutput() bool {
	return true
}

func (this *recipeOutput) WriteHash(hash hash.Hash) {
	hash.Write([]byte("output"))
	this.script.WriteHash(hash)

	if this.always {
		hash.Write([]byte("always"))
	}
	if this.try {
		hash.Write([]byte("try"))
	}
	if this.interpreter != nil {
		this.interpreter.WriteHash(hash)
	}
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
	outpath := createOutDir(sum)

	if _, err := os.Stat(outpath); err == nil {
		if !this.always && !ctx.forceBuild {
			return outpath, nil
		}
		if err = os.RemoveAll(outpath); err != nil {
			return "", err
		}
	}

	workdir, err := os.MkdirTemp(os.TempDir(), "paccat-workdir-******")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workdir) /* do remove the workdir if not needed */

	ctx.Set("out", outpath)
	defer ctx.Unset("out")

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

	return outpath, nil
}

func (this *recipeOutput) GetPosition() position {
	return this.pos
}
