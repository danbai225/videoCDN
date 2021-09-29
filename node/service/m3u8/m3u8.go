package m3u8

import (
	"errors"
	"fmt"
	logs "github.com/danbai225/go-logs"
	m3u8s "github.com/grafov/m3u8"
	"net/url"
	"p00q.cn/video_cdn/comm/utils"
	"p00q.cn/video_cdn/node/config"
	"p00q.cn/video_cdn/node/service/cache"
	downloadServer "p00q.cn/video_cdn/node/service/download"
	"strings"
)

//解析m3u8文件链接
func parseM3U8Url(m3u8Url string, mediaType m3u8s.ListType) (*url.URL, interface{}, error) {
	urlParse, err := url.Parse(m3u8Url)
	if err != nil {
		return nil, nil, err
	}
	resp, err := downloadServer.Get(m3u8Url)
	if err != nil {
		return urlParse, nil, err
	}
	if resp.StatusCode != 200 {
		return urlParse, nil, errors.New(fmt.Sprintf("未能正确获取到资源，状态为%d", resp.StatusCode))
	}
	list, listType, err := m3u8s.DecodeFrom(resp, true)
	if err != nil {
		return urlParse, nil, err
	}
	if listType != mediaType {
		return urlParse, nil, errors.New("不是预期的ListType")
	}
	return urlParse, list, nil
}

//获取加密key
func getKey(urlP *url.URL, playlist *m3u8s.MediaPlaylist) (string, error) {
	if playlist.Key.URI != "" {
		parse, err2 := url.Parse(playlist.Key.URI)
		if err2 != nil {
			return "", err2
		}
		KeyUrl := playlist.Key.URI
		if !parse.IsAbs() {
			KeyUrl = fmt.Sprintf("%s://%s:%s%s", urlP.Scheme, urlP.Host, urlP.Port(), KeyUrl)
		}
		download, err := downloadServer.Download(KeyUrl)
		return string(download), err
	}
	return "", nil
}
func NewTransit(m3u8Url string) (string, error) {
	urlParse, playlist, err := parseM3U8Url(m3u8Url, m3u8s.MASTER)
	if err != nil {
		return "", err
	}
	videoKey := utils.MD5(urlParse.Host + urlParse.Path)
	master := playlist.(*m3u8s.MasterPlaylist)
	for i := range master.Variants {
		variant := master.Variants[i]
		variantUrlParse, err2 := url.Parse(variant.URI)
		if err2 != nil {
			return "", err2
		}
		Url := variant.URI
		if !variantUrlParse.IsAbs() {
			Url = fmt.Sprintf("%s://%s:%s%s", urlParse.Scheme, urlParse.Host, urlParse.Port(), variant.URI)
		}
		mUrlParse, me, err3 := parseM3U8Url(Url, m3u8s.MEDIA)
		if err3 != nil {
			return "", err3
		}
		mediaList := me.(*m3u8s.MediaPlaylist)
		Key, err4 := getKey(mUrlParse, mediaList)
		if err4 != nil {
			return "", err4
		}
		bkurl := variant.URI
		keyUrl := ""
		if Key != "" {
			keyUrl = fmt.Sprintf("/video/%s/list%d.key", videoKey, i)
			if utils.IsRelativeUrl(bkurl) {
				if utils.IsRelativeUrl(variant.URI) {
					bkurl = fmt.Sprintf("%s://%s:%s%s", urlParse.Scheme, urlParse.Host, urlParse.Port(), bkurl)
				} else {
					parse, err5 := url.Parse(variant.URI)
					if err5 != nil {
						return "", err5
					}
					bkurl = fmt.Sprintf("%s://%s:%s%s", parse.Scheme, parse.Host, parse.Port(), bkurl)
				}
			}
			mediaList.Key.URI = keyUrl
			cache.CacheKey(keyUrl, Key)
		}
		for j := range mediaList.Segments {
			segment := mediaList.Segments[j]
			if segment != nil {
				tsUrl := fmt.Sprintf("/video/%s/list%d/%d.ts", videoKey, i, j)
				if utils.IsRelativeUrl(segment.URI) {
					if utils.IsRelativeUrl(variant.URI) {
						segment.URI = fmt.Sprintf("%s://%s:%s%s", urlParse.Scheme, urlParse.Host, urlParse.Port(), segment.URI)
					} else {
						parse, err5 := url.Parse(bkurl)
						if err5 != nil {
							return "", err5
						}
						segment.URI = fmt.Sprintf("%s://%s:%s%s", parse.Scheme, parse.Host, parse.Port(), segment.URI)
					}
				}
				cache.CacheKey(tsUrl, segment.URI)
				cache.CacheKey(segment.URI, tsUrl)
				segment.URI = tsUrl
			}
		}
		mediaList.ResetCache()
		meUrl := fmt.Sprintf("/video/%s/list%d.m3u8", videoKey, i)
		if keyUrl != "" {
			cache.Cache(meUrl, []byte(strings.ReplaceAll(mediaList.Encode().String(), bkurl, keyUrl)))
		}
		master.Variants[i].URI = meUrl
	}
	master.ResetCache()
	data := master.Encode().Bytes()
	logs.Info("\n", string(data))
	maUrl := fmt.Sprintf("/video/%s/index.m3u8", videoKey)
	cache.Cache(maUrl, data)
	return config.GlobalConfig.Host + maUrl, nil
}

// GetResources 获取m3u8资源文件
func GetResources(key string) ([]byte, error) {
	suffix := strings.HasSuffix(key, "key")
	if suffix {
		//加密key 存在内存中
		return []byte(cache.GetCacheKey(key)), nil
	}
	suffix = strings.HasSuffix(key, "ts")
	if suffix {
		key = cache.GetCacheKey(key)
	}
	if cache.HasKey(key) {
		d, err := cache.GetCache(key)
		if err != nil {
			return nil, err
		}
		return d, nil
	} else {
		formUrl, err := cache.CacheFormUrl(key)
		if err != nil {
			return nil, err
		}
		return formUrl, nil
	}
}
