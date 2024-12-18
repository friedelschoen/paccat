package main

import (
	"fmt"
	"log"
	"os"

	"friedelschoen.io/upkg/recipe"
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
	if len(os.Args) != 2 {
		log.Fatal("usage: upkg <recipe>")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	got, err := recipe.ParseRecipe(os.Args[1], file)
	if err != nil {
		log.Fatal(err)
	}

	result, err := got.Build("build", true, nil)
	if err != nil {
		log.Fatal("error while building: ", err)
	}
	fmt.Println(result)

	if err = makeSymlink(result); err != nil {
		log.Print(err)
	}
}
