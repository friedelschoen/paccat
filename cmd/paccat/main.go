package main

import (
	_ "embed"
	"fmt"
	"os"

	"friedelschoen.io/paccat/internal/errors"
	"friedelschoen.io/paccat/internal/parser"
	"friedelschoen.io/paccat/internal/types"
)

//go:embed cat.txt
var logo string

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
	filename := os.Args[1]
	eval, err := parser.ParseFile(filename)
	if err != nil {
		errors.PrintTrace(os.Stdout, err)
		os.Exit(1)
	}

	ctx := types.NewContext(filename)
	value, err := ctx.Evaluate(eval)
	if err != nil {
		errors.PrintTrace(os.Stdout, err)
		os.Exit(1)
	}

	strValue, err := types.CastString(value, ctx)
	if err != nil {
		errors.PrintTrace(os.Stdout, err)
		os.Exit(1)
	}

	fmt.Println(strValue.Content)
}
