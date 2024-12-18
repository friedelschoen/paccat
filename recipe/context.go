package recipe

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type Context struct {
	currentRecipe *Recipe
	attributes    map[string]Buildable

	// used by FunctionCall
	nextAttribute *string
}

func createOutDir(name string) string {
	nametime := fmt.Sprintf("%s-%d", name, time.Now().UnixMilli())
	home, err := os.UserHomeDir()
	if err != nil {
		return path.Join(os.TempDir(), "upkg", nametime)
	}
	return path.Join(home, ".upkg", nametime)
}

func createWorkdir(name string) string {
	nametime := fmt.Sprintf("%s-%d", name, time.Now().UnixMilli())
	home, err := os.UserHomeDir()
	if err != nil {
		return path.Join(os.TempDir(), "upkg", nametime)
	}
	return path.Join(home, ".upkg", nametime)
}

func (this *Context) Get(key string, forceOutput bool) (string, error) {
	value, ok := this.attributes[key]
	if !ok {
		return "", fmt.Errorf("`%s` not found in current context", key)
	}

	if !value.HasOutput() {
		if forceOutput {
			return "", fmt.Errorf("recipe did not produce output")
		}
		return value.Build(this)
	}

	fmt.Println("output!")

	name, err := this.Get("name", false)
	if err != nil {
		return "", err
	}

	outdir := createOutDir(name)
	workdir, err := os.MkdirTemp(os.TempDir(), "upkg-workdir-******")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workdir) /* do remove the workdir if not needed */

	this.attributes["out"] = &RecipeStringLiteral{outdir}

	script, err := value.Build(this)
	if err != nil {
		return "", fmt.Errorf("error while fetching '%s': %v", key, err)
	}

	fmt.Printf("script: '%v'\n", script)

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
