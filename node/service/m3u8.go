package service

import (
	"github.com/gogf/gf/os/gcache"
	"p00q.cn/video_cdn/comm/model"
	"strings"
	"time"
)

var cacheMap = gcache.New()

// GetResources 获取m3u8资源文件
func GetResources(url string) ([]byte, error) {
	split := strings.Split(url, "/")
	if len(split) < 2 {
		return []byte{}, nil
	}
	key := split[2]
	LoadCacheData(key)
	if HasKey(url) {
		d, err := GetCache(url)
		if err != nil {
			return nil, err
		}
		return d, nil
	} else {
		formUrl, err := CacheFormUrl(url)
		if err != nil {
			if strings.Contains(err.Error(), "Client.Timeout") {
				return CacheFormUrl(url)
			}
		}
		return formUrl, nil
	}
}
func GetUrl(urlKey string) string {
	get, err := cacheMap.Get(urlKey)
	if err != nil && get != nil {
		return ""
	}
	return get.(string)
}
func LoadCacheData(videoKey string) {
	v, err := cacheMap.Get(videoKey)
	if err != nil || v == nil {
		updateCache(GetVideoCacheData(videoKey))
		cacheMap.Set(videoKey, true, time.Hour)
	}
}
func updateCache(data []model.Data) {
	for _, datum := range data {
		switch datum.Type {
		case "data":
			Cache(datum.Key, []byte(datum.Data))
		case "url":
			cacheMap.Set(datum.Key, datum.Data, time.Hour)
		}
	}
}
func clearCacheMap() {
	_ = cacheMap.Clear()
}
