package recipe

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"friedelschoen.io/paccat/internal/util"
)

type Context struct {
	workDir         string               // directory of the recipe
	currentRecipe   *Recipe              // current recipe
	scope           map[string]Evaluable // variables and attributes
	importAttribute *string              // used by RecipeImport
}

func createOutDir(name string) string {
	nametime := fmt.Sprintf("%s-%d", name, time.Now().UnixMilli())
	return path.Join(util.GetCachedir(), nametime)
}

func (this *Context) Get(key string, forceOutput bool) (string, error) {
	value, ok := this.scope[key]
	if !ok {
		return "", NoAttributeError
	}

	if !value.HasOutput() {
		if forceOutput {
			return "", NoOutputError
		}
		return value.Eval(this)
	}

	name, err := this.Get("name", false)
	if err != nil {
		return "", err
	}

	outdir := createOutDir(name)
	workdir, err := os.MkdirTemp(os.TempDir(), "paccat-workdir-******")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workdir) /* do remove the workdir if not needed */

	this.scope["out"] = &recipeStringLiteral{outdir}

	script, err := value.Eval(this)
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
