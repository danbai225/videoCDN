package utils

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func Sha256File(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)
	h := sha256.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
