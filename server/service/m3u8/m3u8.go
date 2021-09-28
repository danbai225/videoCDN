package m3u8

import (
	"fmt"
	m3u8s "github.com/grafov/m3u8"
	"net/url"
	downloadServer "p00q.cn/video_cdn/server/service/download"
	"strings"
)

//解析最终host 可能会携带端口
func parseFinalHost(m3u8Url string) (string, error) {
	urlParse, err := url.Parse(m3u8Url)
	if err != nil {
		return "", err
	}
	resp, err := downloadServer.Get(m3u8Url)
	if err != nil {
		return "", err
	}
	from, listType, err := m3u8s.DecodeFrom(resp, true)
	if err != nil {
		return "", err
	}
	if listType == m3u8s.MASTER {
		playlist := from.(*m3u8s.MasterPlaylist)
		if len(playlist.Variants) > 0 {
			uri := playlist.Variants[0].URI
			if strings.Contains(uri, "://") {
				return uri, nil
			}
			if urlParse.Port() == "" {
				return parseFinalHost(fmt.Sprintf("%s://%s%s", urlParse.Scheme, urlParse.Host, uri))
			}
			return parseFinalHost(fmt.Sprintf("%s://%s:%s%s", urlParse.Scheme, urlParse.Host, urlParse.Port(), uri))
		}
	}
	playlist := from.(*m3u8s.MediaPlaylist)
	if len(playlist.Segments) > 0 {
		if len(playlist.Segments) > 0 {
			segment := playlist.Segments[0]
			uri := segment.URI
			if strings.Contains(uri, "://") {
				parse, err := url.Parse(uri)
				if err == nil {
					println(parse.Host, parse.Port())
					return parse.Host, nil
				}
			}
		}
	}
	return urlParse.Host, nil
}
func ParseFinalHost(m3u8 string) (string, error) {
	host, err := parseFinalHost(m3u8)
	split := strings.Split(host, ":")
	if len(split) > 0 {
		return split[0], err
	}
	return "", err
}
