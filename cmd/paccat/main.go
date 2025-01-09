package main

import (
	_ "embed"
	"fmt"
	"os"

	"friedelschoen.io/paccat/internal/ast"
	"friedelschoen.io/paccat/internal/errors"
	"friedelschoen.io/paccat/internal/parser"
	"friedelschoen.io/paccat/internal/types"
)

//go:embed cat.txt
var logo string

//go:embed help.txt
var helpmsg string

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
	printast := false
	printsource := false
	makeresult := false
	i := 0
argloop:
	for i = 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--help", "-h":
			fmt.Print(helpmsg)
			os.Exit(0)
		case "--ast", "-t":
			printast = true
		case "--source", "-s":
			printsource = true
		case "--result":
			makeresult = true
		default:
			if os.Args[i][0] == '-' {
				fmt.Fprintf(os.Stderr, "error: unknown option '%s'\n%s", os.Args[i], helpmsg)
				os.Exit(1)
			}
			break argloop
		}
	}
	if i != len(os.Args)-1 {
		fmt.Fprint(os.Stderr, helpmsg)
		os.Exit(1)
	}

	filename := os.Args[i]
	eval, err := parser.ParseFile(filename)
	if err != nil {
		errors.PrintTrace(os.Stdout, err)
		os.Exit(1)
	}

	if printast {
		ast.PrintTree(os.Stdout, eval, 0)
	}

	ctx := types.Scope{}
	value, err := ctx.Evaluate(eval)
	if err != nil {
		errors.PrintTrace(os.Stdout, err)
		os.Exit(1)
	}

	fmt.Println(value.Content)

	if printsource {
		for _, ss := range value.FlatSources() {
			fmt.Printf("%d-%d: %s\n", ss.Start, ss.Start+ss.Len, ss.Value.Node.Name())
		}
	}

	if makeresult {
		fmt.Println(value.Content)
		if makeresult {
			path := value.Content
			if _, err := os.Lstat(path); err != nil {
				fmt.Fprintf(os.Stderr, "error: unable to stat result: %v\n", err)
				os.Exit(1)
			}
			makeSymlink(path)
		}
	}
}
