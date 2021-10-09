package utils

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

func HostAddPath(url *url.URL, path string) string {
	base := filepath.Dir(url.Path)
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	if base != "/" && !strings.HasPrefix(path, base) {
		path = base + path
	}
	if url.Port() == "" {
		return fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, path)
	}
	return fmt.Sprintf("%s://%s:%s%s", url.Scheme, url.Host, url.Port(), path)
}
