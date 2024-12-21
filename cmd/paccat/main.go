package main

import (
	_ "embed"
	"flag"
	"fmt"
	"hash/maphash"
	"log"
	"os"
	"path"

	"friedelschoen.io/paccat/internal/recipe"
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
	install := flag.Bool("install", false, "Build the package")
	evaluate := flag.String("evaluate", "", "Evaluate attribute")
	dohash := flag.Bool("hash", false, "Hash Recipe and return")
	noResult := flag.Bool("no-result", false, "Don't expect a path")

	flag.BoolFunc("help", "prints help-message", func(string) error {
		fmt.Print(logo)
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
	if err != nil {
		log.Fatal(err)
	}

	if *install {
		ctx.BuildPackage()
		if err != nil {
			log.Fatal("error while building: ", err)
		}
	} else if *evaluate != "" {
		result, err := ctx.Get(*evaluate, *noResult)

		if err != nil {
			log.Fatal("error while building: ", err)
		}
		fmt.Println(result)

		if !*noResult {
			if err = makeSymlink(result); err != nil {
				log.Print(err)
			}
		}
	} else if *dohash {
		hash := maphash.Hash{}
		ast.(*recipe.Recipe).WriteHash(hash)
		sum := hash.Sum64()
		fmt.Printf("%016x", sum)
	} else {
		fmt.Printf("no operation\n")
	}
}
