package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"friedelschoen.io/upkg/internal/recipe"
)

func makeSymlink(result string) error {
	// Check if the file or directory exists
	info, err := os.Lstat("result")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to stat ./result: %v", err)
	}

	if err == nil {
		// Check if the existing path is a symlink
		if info.Mode()&os.ModeSymlink == 0 { // Path exists and is not a symlink - throw an error
			return fmt.Errorf("path ./result exists and is not a symlink")
		}

		// Path is a symlink, remove it
		if err := os.Remove("result"); err != nil {
			return fmt.Errorf("failed to remove symlink ./result: %v", err)
		}
	}

	return os.Symlink(result, "result")
}

func main() {
	evaluate := flag.Bool("evaluate", false, "Evaluate any value")
	attribute := flag.String("attribute", "build", "Attribute to evaluate")
	noResult := flag.Bool("no-result", false, "Don't symlink to ./result")
	flag.BoolFunc("help", "prints help-message", func(string) error {
		flag.Usage()
		os.Exit(0)
		return nil
	})
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filepath := flag.Arg(0)

	ast, err := recipe.ParseFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := ast.(*recipe.Recipe).NewContext(path.Dir(filepath), nil)
	result, err := ctx.Get(*attribute, !*evaluate)

	if err != nil {
		log.Fatal("error while building: ", err)
	}
	fmt.Println(result)

	if !*evaluate && !*noResult {
		if err = makeSymlink(result); err != nil {
			log.Print(err)
		}
	}
}
