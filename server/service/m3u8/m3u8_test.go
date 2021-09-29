package m3u8

import (
	"p00q.cn/video_cdn/server/global"
	"testing"
)

func TestParseFinalHost(t *testing.T) {
	println(ParseFinalHost("https://vod1.bdzybf1.com/20201015/bwcEwYbn/index.m3u8"))
}
func TestParseM3U8AndCacheTheURL(t *testing.T) {
	global.InitDB()
	//urlt:="https://vod1.bdzybf1.com/20201015/bwcEwYbn/index.m3u8"
	//urlt:="https://t.wdubo.com/20210830/cWNlCZNT/index.m3u8"
	urlt := "http://creakvipok9baiduyd.czmhgz.cn/03/veto-free.guy.2021.1080p.bluray.x264/pl.m3u8"
	url, err := ParseM3U8AndCacheTheURL(urlt)
	if err != nil {
		t.Fatal(err.Error())
	}
	println(url)
}
func TestCache(t *testing.T) {
	//println(CacheM3u8("https://vod1.bdzybf1.com/20201015/bwcEwYbn/index.m3u8"))
}
