package main

import (
	"flag"
	"log"
	"os"
	"path"

	"friedelschoen.io/upkg/internal/recipe"
)

func main() {
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
	_, err = ctx.BuildPackage()
	if err != nil {
		log.Fatal(err)
	}
}
