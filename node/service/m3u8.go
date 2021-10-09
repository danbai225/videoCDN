package service

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gcache"
	"math/rand"
	"net/url"
	"p00q.cn/video_cdn/comm/model"
	"path"
	"strings"
	"sync"
	"time"
)

var cacheMap = gcache.New()

type Caches struct {
	HeadId int64
	Val    string
}

// GetResources 获取m3u8资源文件
func GetResources(url string) ([]byte, error) {
	key := getVideoKey(url)
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
func getVideoKey(url string) string {
	split := strings.Split(url, "/")
	if len(split) < 2 {
		return ""
	}
	return split[2]
}

var lockMap = gmap.New(true)

func LoadCacheData(videoKey string) {
	lock := lockMap.GetOrSetFuncLock(videoKey, func() interface{} {
		mutex := &sync.Mutex{}
		return mutex
	})
	lock.(*sync.Mutex).Lock()
	defer func() {
		lock.(*sync.Mutex).Unlock()
		lockMap.Remove(videoKey)
	}()
	v, err := cacheMap.Get(videoKey)
	if err != nil || v == nil {
		updateCache(GetVideoCacheData(videoKey))
		cacheMap.Set(videoKey, true, time.Hour)
	}
}
func getShortKey(key string) string {
	if len(key) <= 40 {
		return ""
	}
	return key[40:]
}
func updateCache(data []model.Data) {
	urls := make([]string, 0)
	head := ""
	id := int64(0)
	val := ""
	for _, datum := range data {
		///video/024fccb26431faf5c8b33cbb7a8989c2/list0/1007.ts
		switch datum.Type {
		case "data":
			Cache(datum.Key, []byte(datum.Data))
		case "url":
			if head == "" {
				urls = append(urls, datum.Data)
				if len(urls) == 3 {
					head = extractCommonHead(urls)
					if head == "" {
						head = "host"
					} else {
						id = addHead(head)
					}
				}
			}
			//去主机段组成公共部分
			if head == "host" || head == "" {
				hostHead := getHostHead(datum.Data)
				val = strings.ReplaceAll(datum.Data, hostHead, "")
				headIdKey := fmt.Sprintf("hedeId-%s", hostHead)
				has, _ := cacheMap.Contains(headIdKey)
				if has {
					get, err := cacheMap.Get(headIdKey)
					if get != nil && err == nil {
						id = get.(int64)
					}
				} else {
					id = rand.Int63()
					key := fmt.Sprintf("head-%d", id)
					hasId := true
					for hasId {
						hasId, _ = cacheMap.Contains(key)
						if hasId {
							id++
							key = fmt.Sprintf("head-%d", id)
						}
					}
					cacheMap.Set(headIdKey, id, 0)
					cacheMap.Set(key, hostHead, 0)
				}
			} else {
				val = strings.ReplaceAll(datum.Data, head, "")
			}
			caches := Caches{
				HeadId: id,
				Val:    val,
			}
			cacheMap.Set(getShortKey(datum.Key), caches, time.Hour)
		}
	}
}
func GetCacheMapData(key string) string {
	get, err := cacheMap.Get(getShortKey(key))
	if err == nil && get != nil {
		caches := get.(Caches)
		headKey := fmt.Sprintf("head-%d", caches.HeadId)
		head, err := cacheMap.Get(headKey)
		if err == nil && head != nil {
			return head.(string) + caches.Val
		}
	}
	return ""
}
func getHostHead(urlString string) string {
	parse, _ := url.Parse(urlString)
	return fmt.Sprintf("%s://%s", parse.Scheme, parse.Host)
}
func addHead(head string) int64 {
	id := rand.Int63()
	key := fmt.Sprintf("head-%d", id)
	has := true
	for has {
		has, _ = cacheMap.Contains(key)

		if has {
			id++
			key = fmt.Sprintf("head-%d", id)
		}
	}
	cacheMap.Set(key, head, 0)
	return id
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
func clearCacheMap() {
	_ = cacheMap.Clear()
}
func GetUrl(urlKey string) string {
	return GetCacheMapData(urlKey)
}
