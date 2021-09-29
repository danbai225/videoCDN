package utils

import (
	"fmt"
	"net/url"
)

func HostAddPath(url *url.URL, path string) string {
	if url.Port() == "" {
		return fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, path)
	}
	return fmt.Sprintf("%s://%s:%s%s", url.Scheme, url.Host, url.Port(), path)
}
