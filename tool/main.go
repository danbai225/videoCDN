package main

import (
	"net/url"
	"path"
)

func main() {
	//clients(1, true, true)
	s := extractCommonHead([]string{"https://ts8.hhmm0.com:9999/20211004/4gSsH6QF/1000kb/hls/bmwH7zUG.ts?a=21", "https://ts8.hhmm0.com:9999/20211004/4gSsH6QF/1000kb/hls/q5HQXOPR.ts"})
	println(s)
}
func extractCommonHead(urls []string) string {
	if len(urls) < 2 {
		return ""
	}
	ps := make([]*url.URL, 0)
	for _, s := range urls {
		parse, _ := url.Parse(s)
		ps = append(ps, parse)
	}
	Scheme := ps[0].Scheme
	Host := ps[0].Host
	Path := path.Dir(ps[0].Path)
	pathFlg := false
	for i := 1; i < len(ps); i++ {
		if ps[i].Scheme != Scheme {
			return ""
		}
		if ps[i].Host != Host {
			return ""
		}
		if path.Dir(ps[i].Path) != Path {
			pathFlg = true
		}
	}
	if pathFlg {
		return Scheme + "://" + Host
	}
	return Scheme + "://" + Host + Path
}
