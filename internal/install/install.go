package install

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type PackageDatabase struct {
	Prefix string
}

// install installs a package by creating directories and symlinking files. It tracks all created entries.
func (db *PackageDatabase) Install(pkgname, pathname string) error {
	os.MkdirAll(db.Prefix, 0755)

	dbpath := path.Join(db.Prefix, "paccat.index")
	file, err := os.OpenFile(dbpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	fmt.Printf("installing %s\n", pkgname)

	err = filepath.Walk(pathname, func(currentPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(pathname, currentPath)
		if err != nil {
			return err
		}

		targetPath := path.Join(db.Prefix, relPath)
		if info.IsDir() {
			csvWriter.Write([]string{pkgname, "link", targetPath})
			if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
				return err
			}
		} else {
			csvWriter.Write([]string{pkgname, "dir", targetPath})
			if err := os.Symlink(currentPath, targetPath); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// remove removes all files and directories for a package but keeps non-empty directories.
func (db *PackageDatabase) Remove(pkgname string) error {
	oldpath := path.Join(db.Prefix, "paccat.index")
	oldfile, err := os.Open(oldpath)
	if err != nil {
		return err
	}
	defer oldfile.Close()

	newpath := path.Join(db.Prefix, "paccat.index.new")
	newfile, err := os.OpenFile(newpath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer newfile.Close()
	defer os.Rename(newpath, oldpath)

	reader := csv.NewReader(oldfile)
	writer := csv.NewWriter(newfile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		name := record[0]
		if name != pkgname {
			writer.Write(record)
			continue
		}
		fmt.Printf("removing %s\n", name)

		target := record[2]

		if err = os.Remove(target); err != nil {
			fmt.Fprintf(os.Stderr, "unable to remove %s %s: %v\n", record[1], target, err)
		}
	}

	return nil
}
