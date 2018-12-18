package utils

import (
	"os"
	"path"
	"path/filepath"
)

func FindUniqueFileNameInPath(source string, targetDir string) (out string, unique string) {
	var filename = filepath.Base(source)
	var extension = filepath.Ext(source)
	var prefixedFilename = filename[0 : len(filename)-len(extension)]
	for {
		unique = RandomString(32)
		out = path.Join(targetDir, prefixedFilename+"_"+unique+extension)

		if _, err := os.Stat(out); os.IsNotExist(err) {
			break
		}
	}

	return
}
