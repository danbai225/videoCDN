package m3u8

import (
	"errors"
	"fmt"
	m3u8s "github.com/grafov/m3u8"
	"net/url"
	"p00q.cn/video_cdn/comm/model"
	"p00q.cn/video_cdn/comm/utils"
	"p00q.cn/video_cdn/server/global"
	downloadServer "p00q.cn/video_cdn/server/service/download"
	nodeServer "p00q.cn/video_cdn/server/service/node"
	"strings"
	"time"
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
				parse, err2 := url.Parse(uri)
				if err2 == nil {
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

//解析m3u8文件链接
func getM3U8UrlContent(m3u8Url string) (*url.URL, interface{}, m3u8s.ListType, error) {
	urlParse, err := url.Parse(m3u8Url)
	if err != nil {
		return nil, nil, 0, err
	}
	resp, err := downloadServer.Get(m3u8Url)
	if err != nil {
		return urlParse, nil, 0, err
	}
	if resp.StatusCode != 200 {
		return urlParse, nil, 0, errors.New(fmt.Sprintf("未能正确获取到资源，状态为%d", resp.StatusCode))
	}
	list, listType, err := m3u8s.DecodeFrom(resp, true)
	if err != nil {
		return urlParse, nil, 0, err
	}
	return urlParse, list, listType, nil
}

// ParseM3U8AndCacheTheURL 对m3u8播放列表进行url进行替换
func ParseM3U8AndCacheTheURL(m3u8 string) (string, error) {
	var err error
	urlP, list, listType, err := getM3U8UrlContent(m3u8)
	if err != nil {
		return "", err
	}
	videoKey := utils.MD5(urlP.Host + urlP.Path)
	defer func() {
		if err != nil {
			global.MySQL.Where("video_key=?", videoKey).Delete(&model.Data{})
		}
	}()
	if listType == m3u8s.MASTER {
		playlist := list.(*m3u8s.MasterPlaylist)
		variants := playlist.Variants
		if len(variants) == 0 {
			return "", errors.New("未获取到播放列表")
		}
		for i, variant := range variants {
			if !strings.Contains(variant.URI, "://") {
				variant.URI = utils.HostAddPath(urlP, variant.URI)
			}
			_, err2 := parseMediaM3U8AndCacheTheURL(videoKey, variant.URI, i)
			if err2 != nil {
				return "", err2
			}
			variant.URI = fmt.Sprintf("/video/%s/list%d.m3u8", videoKey, i)
		}
		cacheUrl := fmt.Sprintf("/video/%s/index.m3u8", videoKey)
		playlist.ResetCache()
		return cacheUrl, global.MySQL.Model(&model.Data{}).Create(&model.Data{
			Key:      cacheUrl,
			VideoKey: videoKey,
			Type:     "data",
			Data:     playlist.Encode().String(),
		}).Error
	}
	return parseMediaM3U8AndCacheTheURL(videoKey, m3u8, 0)
}
func parseMediaM3U8AndCacheTheURL(videoKey string, m3u8 string, index int) (string, error) {
	urlP, list, listType, err := getM3U8UrlContent(m3u8)
	if err != nil {
		return "", err
	}
	if listType != m3u8s.MEDIA {
		return "", errors.New("类型错误 listType!=m3u8s.MEDIA")
	}
	mediaList := list.(*m3u8s.MediaPlaylist)
	segments := mediaList.Segments
	if len(segments) == 0 {
		return "", errors.New("未获取到播放列表资源")
	}
	ds := make([]model.Data, 0)
	//存在加密
	keyUrl := mediaList.Key.URI
	if mediaList.Key != nil && mediaList.Key.URI != "" {

		if !strings.Contains(keyUrl, "://") {
			keyUrl = utils.HostAddPath(urlP, keyUrl)
		}
		keyBytes, err2 := downloadServer.Download(keyUrl)
		if err2 != nil {
			return "", err2
		}
		mediaList.Key.URI = fmt.Sprintf("/video/%s/list%d.key", videoKey, index)
		ds = append(ds, model.Data{
			Key:      mediaList.Key.URI,
			VideoKey: videoKey,
			Type:     "data",
			Data:     string(keyBytes),
		})
	}
	//缓存ts url
	for i, segment := range segments {
		if segment == nil {
			continue
		}
		if !strings.Contains(segment.URI, "://") {
			segment.URI = utils.HostAddPath(urlP, segment.URI)
		}
		cacheUrl := fmt.Sprintf("/video/%s/list%d/%d.ts", videoKey, index, i)
		ds = append(ds, model.Data{
			Key:      cacheUrl,
			VideoKey: videoKey,
			Type:     "url",
			Data:     segment.URI,
		})
		segment.URI = cacheUrl
	}
	mediaList.ResetCache()
	cacheUrl := fmt.Sprintf("/video/%s/list%d.m3u8", videoKey, index)
	s := mediaList.Encode().String()
	s = strings.ReplaceAll(s, keyUrl, mediaList.Key.URI)
	ds = append(ds, model.Data{
		Key:      cacheUrl,
		VideoKey: videoKey,
		Type:     "data",
		Data:     s,
	})
	return cacheUrl, global.MySQL.Create(&ds).Error
}

// CacheM3u8 获取缓存m3u8地址
func CacheM3u8(m3u8 string) (string, error) {
	now := time.Now()
	cache := model.Cache{}
	global.MySQL.Model(&model.Cache{}).Where("url=?", m3u8).Take(&cache)
	if cache.ID != 0 {
		return cache.NodeUrl, nil
	}
	host, err := ParseFinalHost(m3u8)
	if err != nil {
		return "", err
	}
	global.Logs.Info(0, fmt.Sprintf("%.2f", time.Now().Sub(now).Seconds()))
	nodeServer.DelayTest(host)
	global.Logs.Info(1, fmt.Sprintf("%.2f", time.Now().Sub(now).Seconds()))
	thePath, err := ParseM3U8AndCacheTheURL(m3u8)
	if err != nil {
		return "", err
	}
	global.Logs.Info(2, fmt.Sprintf("%.2f", time.Now().Sub(now).Seconds()))
	delay := model.Delay{}
	err = global.MySQL.Model(&model.Delay{}).Where("host=?", host).Order("val ASC").First(&delay).Error
	var node model.Node
	if delay.Val > 0 && err != nil {
		node = nodeServer.GetNodeInfoByIP(delay.NodeIP)
	} else {
		node = nodeServer.AssignANodeWithTheLeastLoad()
		if node.ID == 0 {
			return "", errors.New("没有在线的节点")
		}
	}
	protocol := "http"
	if node.Https {
		protocol = "https"
	}
	hostUrl := node.IP
	if node.Domain != "" {
		hostUrl = node.Domain
	}
	global.Logs.Info(3, fmt.Sprintf("%.2f", time.Now().Sub(now).Seconds()))
	cacheUrl := fmt.Sprintf("%s://%s:%d%s", protocol, hostUrl, node.Port, thePath)
	urlP, _ := url.Parse(m3u8)
	videoKey := utils.MD5(urlP.Host + urlP.Path)
	global.MySQL.Create(&model.Cache{
		NodeIP:   node.IP,
		Host:     host,
		NodeUrl:  cacheUrl,
		Url:      m3u8,
		Visits:   0,
		Valid:    true,
		Flow:     0,
		VideoKey: videoKey,
	})
	global.Logs.Info(4, fmt.Sprintf("%.2f", time.Now().Sub(now).Seconds()))
	go nodeServer.NewCacheData(videoKey, node.IP)
	return cacheUrl, nil
}
