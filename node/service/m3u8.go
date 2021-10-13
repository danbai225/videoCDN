package service

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/os/gcache"
	"net/url"
	"p00q.cn/video_cdn/comm/model"
	"path"
	"strconv"
	"strings"
	"time"
)

var cacheMap = gcache.New()

type Caches struct {
	HeadId byte
	Val    string
}
type VideoUrlCache struct {
	heads map[byte]string
	url   map[string]Caches
}

var downloadSet = gset.New(true)

// GetResources 获取m3u8资源文件
func GetResources(url string) ([]byte, error) {
	if !downloadSet.AddIfNotExist(url) {
		for !downloadSet.AddIfNotExist(url) {
			time.Sleep(time.Millisecond * 100)
		}
	}
	defer downloadSet.Remove(url)
	key := getVideoKey(url)
	LoadCacheData(key)
	if HasKey(url) {
		d, err := GetCache(url)
		if err != nil {
			return nil, err
		}
		return d, nil
	} else {
		last10key := getTheLast10key(url)
		//logs.Info(len(last10key))
		for _, s := range last10key {
			cacheUrl := s
			if downloadSet.AddIfNotExist(cacheUrl) {
				//logs.Info("提前缓存",cacheUrl)
				go func() {
					CacheFormUrl(cacheUrl)
					downloadSet.Remove(cacheUrl)
				}()
			}
		}
		formUrl, err := CacheFormUrl(url)
		if err != nil {
			if strings.Contains(err.Error(), "Client.Timeout") {
				return CacheFormUrl(url)
			}
			return formUrl, err
		}
		return formUrl, nil
	}
}
func getTheLast10key(url string) []string {
	key := getVideoKey(url)
	k := fmt.Sprintf("url-%s", key)
	get, err := cacheMap.Get(k)
	arr := make([]string, 0)
	if err == nil && get != nil {
		vc := get.(VideoUrlCache)
		split := strings.Split(url, "/")
		ts := split[len(split)-1]
		numTs := strings.Split(ts, ".")
		num, _ := strconv.Atoi(numTs[0])
		if num%10 == 0 {
			tUrl := strings.Join(split[:len(split)-1], "/")
			for i := 1; i < 10; i++ {
				newUrl := tUrl + fmt.Sprintf("/%d.ts", i+num)
				_, has := vc.url[getShortKey(newUrl)]
				if has {
					arr = append(arr, newUrl)
				}
			}
		}
	}
	return arr
}
func getVideoKey(url string) string {
	split := strings.Split(url, "/")
	if len(split) < 2 {
		return ""
	}
	return split[2]
}

func LoadCacheData(videoKey string) {
	if !downloadSet.AddIfNotExist(videoKey) {
		for !downloadSet.AddIfNotExist(videoKey) {
			time.Sleep(time.Millisecond * 100)
		}
	}
	defer downloadSet.Remove(videoKey)
	v, err := cacheMap.Get(fmt.Sprintf("url-%s", videoKey))
	if err != nil || v == nil {
		updateCache(GetVideoCacheData(videoKey))
	}
}
func getShortKey(key string) string {
	if len(key) <= 40 {
		return ""
	}
	return key[40:]
}
func updateCache(data []model.Data) {
	if len(data) == 0 {
		return
	}
	heades := make(map[string]byte)
	cache := VideoUrlCache{
		heads: make(map[byte]string),
		url:   make(map[string]Caches),
	}
	startId := byte(0)
	urls := make([]string, 0)
	head := ""
	id := byte(0)
	val := ""
	videoKey := ""
	for _, datum := range data {
		///video/024fccb26431faf5c8b33cbb7a8989c2/list0/1007.ts
		switch datum.Type {
		case "data":
			videoKey = datum.VideoKey
			Cache(datum.Key, []byte(datum.Data))
		case "url":
			urls = append(urls, datum.Data)
			if len(urls) == 3 {
				head = extractCommonHead(urls)
				vid, has := heades[head]
				if !has {
					cache.heads[startId] = head
					heades[head] = startId
					id = startId
					startId++
				} else {
					id = vid
				}
				urls = make([]string, 0)
			}
			///#https://video.dious.cc/20200617/aAUCQ5Hf/index.m3u8
			//去主机段组成公共部分
			if head == "" || !strings.Contains(datum.Data, head) {
				has := false
				max := 0
				for k, _ := range heades {
					if strings.Contains(datum.Data, k) {
						if max < len(k) {
							has = true
							head = k
						}
					}
				}
				if !has {
					hostHead := getHostHead(datum.Data)
					head = hostHead
					vid, has2 := heades[hostHead]
					if !has2 {
						cache.heads[startId] = head
						heades[head] = startId
						id = startId
						startId++
					} else {
						id = vid
					}
				}
			}
			val = strings.ReplaceAll(datum.Data, head, "")
			id = heades[head]
			caches := Caches{
				HeadId: id,
				Val:    val,
			}
			cache.url[getShortKey(datum.Key)] = caches
		}
	}
	logs.Info("接收缓存url完毕，", len(cache.url), len(data))
	cacheMap.Set(fmt.Sprintf("url-%s", videoKey), cache, time.Minute*10)
}
func GetCacheMapData(key string) string {
	k := fmt.Sprintf("url-%s", getVideoKey(key))
	get, err := cacheMap.Get(k)
	if err == nil && get != nil {
		vc := get.(VideoUrlCache)
		v, has := vc.url[getShortKey(key)]
		if has {
			cacheMap.UpdateExpire(k, time.Minute*6)
			return vc.heads[v.HeadId] + v.Val
		}
	}
	return ""
}

func getHostHead(urlString string) string {
	parse, _ := url.Parse(urlString)
	return fmt.Sprintf("%s://%s", parse.Scheme, parse.Host)
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
	//logs.Info(urlKey,GetCacheMapData(urlKey))
	data := GetCacheMapData(urlKey)
	if data == "" {
		LoadCacheData(getVideoKey(urlKey))
		return GetCacheMapData(urlKey)
	}
	return GetCacheMapData(urlKey)
}
