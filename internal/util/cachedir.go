package util

import (
	"os"
	"path"
)

func GetCachedir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path.Join(os.TempDir(), "paccat")
	}
	return path.Join(home, ".paccat")
}
