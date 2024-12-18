package recipe

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"friedelschoen.io/upkg/internal/util"
)

type Context struct {
	directory     string
	currentRecipe *Recipe
	attributes    map[string]Buildable
	nextAttribute *string
	building      bool
}

func createOutDir(name string) string {
	nametime := fmt.Sprintf("%s-%d", name, time.Now().UnixMilli())
	return path.Join(util.GetCachedir(), nametime)
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
		return "", NoAttributeError
	}

	if !value.HasOutput() {
		if forceOutput {
			return "", NoOutputError
		}
		return value.Build(this)
	}

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
		return "", err
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

func installPath(pathname string) error {
	target := "target"

	return filepath.Walk(pathname, func(currentPath string, info fs.FileInfo, err error) error {
		relPath, err := filepath.Rel(pathname, currentPath)
		if err != nil {
			return err
		}

		targetPath := path.Join(target, relPath)

		if info.IsDir() {
			fmt.Printf("mkdir %s\n", targetPath)
			os.Mkdir(targetPath, info.Mode())
		} else {
			fmt.Printf("symlink %s -> %s\n", currentPath, targetPath)
			os.Symlink(currentPath, targetPath)
		}
		return nil
	})
}

func (this *Context) BuildPackage() (string, error) {
	this.building = true
	defer func() {
		this.building = false
	}()

	buildDepends, err := this.Get("build_depends", false)
	if err != nil && err != NoAttributeError {
		return "", err
	}

	if err == nil {
		for _, dep := range strings.Split(buildDepends, " ") {
			err := installPath(dep)
			if err != nil {
				return "", err
			}
		}
	}

	result, err := this.Get("build", true)
	if err != nil {
		return "", err
	}

	err = installPath(result)
	if err != nil {
		return "", err
	}

	runDepends, err := this.Get("depends", false)
	if err != nil && err != NoAttributeError {
		return "", err
	}

	if err == nil {
		for _, dep := range strings.Split(runDepends, " ") {
			err := installPath(dep)
			if err != nil {
				return "", err
			}
		}
	}

	return result, nil
}
