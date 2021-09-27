package utils

import (
	"path/filepath"
	"strings"
)

func IsAbsPath(path string) bool {
	return filepath.IsAbs(path)
}
func Abs(path string) (string, error) {
	return filepath.Abs(path)
}
func IsRelativeUrl(url string) bool {
	if strings.HasPrefix(url, "http") {
		return false
	}
	return true
}
