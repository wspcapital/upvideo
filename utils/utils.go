package utils

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generate alpha-numeric random string n-length
func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(b)
}

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
