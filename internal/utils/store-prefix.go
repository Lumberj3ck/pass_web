package utils

import (
	"os"
	"path/filepath"
)

func GetStorePrefix() string {
	prefix := os.Getenv("PREFIX")
	if prefix == "" {
		prefix = filepath.Join(os.Getenv("HOME"), ".password-store/")
	}
	return prefix
}
