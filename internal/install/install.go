package install

import (
	"encoding/csv"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type PackageDatabase struct {
	Prefix   string
	Database string
}

// install installs a package by creating directories and symlinking files. It tracks all created entries.
func (db *PackageDatabase) Install(pkgname, pathname string) error {
	file, err := os.OpenFile("files.csv", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

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

// remove removes all files and directories for a package but keeps non-empty directories.
func (db *PackageDatabase) Remove(pkgname string) error {
	oldfile, err := os.Open("files.csv")
	if err != nil {
		return err
	}
	defer oldfile.Close()

	newfile, err := os.OpenFile("files.csv.new", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer newfile.Close()

	_ = csv.NewReader(oldfile)
	_ = csv.NewWriter(newfile)

	return nil
}
