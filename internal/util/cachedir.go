package util

import (
	"log"
	"os"
	"path"
)

func GetCachedir() (dist string) {
	if home, err := os.UserHomeDir(); err == nil {
		dist = path.Join(home, ".paccat/store")
	} else {
		dist = path.Join(os.TempDir(), "paccat")
	}

	if err := os.MkdirAll(dist, 0777); err != nil {
		log.Panic(err)
	}
	return dist
}
