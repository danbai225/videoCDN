package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	logs "github.com/danbai225/go-logs"
	m3u8s "github.com/grafov/m3u8"
	"github.com/xxjwxc/gowp/workpool"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getList(pg int) (List, error) {
	logs.Info("开始获取第", pg, "页")
	get, err := http.Get(fmt.Sprintf("https://api.apibdzy.com/api.php/provide/vod/?ac=list&pg=%d", pg))
	if err != nil {
		return List{}, err
	}
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return List{}, err
	}
	var list List
	json.Unmarshal(all, &list)
	return list, nil
}
func getInfo(id int) (Info, error) {
	get, err := http.Get(fmt.Sprintf("https://api.apibdzy.com/api.php/provide/vod/?ac=detail&ids=%d", id))
	if err != nil {
		return Info{}, err
	}
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return Info{}, err
	}
	var info Info
	json.Unmarshal(all, &info)
	return info, nil
}

type List struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Page      string `json:"page"`
	Pagecount int    `json:"pagecount"`
	Limit     string `json:"limit"`
	Total     int    `json:"total"`
	List      []struct {
		VodId       int    `json:"vod_id"`
		VodName     string `json:"vod_name"`
		TypeId      int    `json:"type_id"`
		TypeName    string `json:"type_name"`
		VodEn       string `json:"vod_en"`
		VodTime     string `json:"vod_time"`
		VodRemarks  string `json:"vod_remarks"`
		VodPlayFrom string `json:"vod_play_from"`
	} `json:"list"`
	Class []struct {
		TypeId   int    `json:"type_id"`
		TypeName string `json:"type_name"`
	} `json:"class"`
}
type Info struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Page      int    `json:"page"`
	Pagecount int    `json:"pagecount"`
	Limit     string `json:"limit"`
	Total     int    `json:"total"`
	List      []struct {
		VodId            int         `json:"vod_id"`
		TypeId           int         `json:"type_id"`
		TypeId1          int         `json:"type_id_1"`
		GroupId          int         `json:"group_id"`
		VodName          string      `json:"vod_name"`
		VodSub           string      `json:"vod_sub"`
		VodEn            string      `json:"vod_en"`
		VodStatus        int         `json:"vod_status"`
		VodLetter        string      `json:"vod_letter"`
		VodColor         string      `json:"vod_color"`
		VodTag           string      `json:"vod_tag"`
		VodClass         string      `json:"vod_class"`
		VodPic           string      `json:"vod_pic"`
		VodPicThumb      string      `json:"vod_pic_thumb"`
		VodPicSlide      string      `json:"vod_pic_slide"`
		VodActor         string      `json:"vod_actor"`
		VodDirector      string      `json:"vod_director"`
		VodWriter        string      `json:"vod_writer"`
		VodBehind        string      `json:"vod_behind"`
		VodBlurb         string      `json:"vod_blurb"`
		VodRemarks       string      `json:"vod_remarks"`
		VodPubdate       string      `json:"vod_pubdate"`
		VodTotal         int         `json:"vod_total"`
		VodSerial        string      `json:"vod_serial"`
		VodTv            string      `json:"vod_tv"`
		VodWeekday       string      `json:"vod_weekday"`
		VodArea          string      `json:"vod_area"`
		VodLang          string      `json:"vod_lang"`
		VodYear          string      `json:"vod_year"`
		VodVersion       string      `json:"vod_version"`
		VodState         string      `json:"vod_state"`
		VodAuthor        string      `json:"vod_author"`
		VodJumpurl       string      `json:"vod_jumpurl"`
		VodTpl           string      `json:"vod_tpl"`
		VodTplPlay       string      `json:"vod_tpl_play"`
		VodTplDown       string      `json:"vod_tpl_down"`
		VodIsend         int         `json:"vod_isend"`
		VodLock          int         `json:"vod_lock"`
		VodLevel         int         `json:"vod_level"`
		VodCopyright     int         `json:"vod_copyright"`
		VodPoints        int         `json:"vod_points"`
		VodPointsPlay    int         `json:"vod_points_play"`
		VodPointsDown    int         `json:"vod_points_down"`
		VodHits          int         `json:"vod_hits"`
		VodHitsDay       int         `json:"vod_hits_day"`
		VodHitsWeek      int         `json:"vod_hits_week"`
		VodHitsMonth     int         `json:"vod_hits_month"`
		VodDuration      string      `json:"vod_duration"`
		VodUp            int         `json:"vod_up"`
		VodDown          int         `json:"vod_down"`
		VodScore         string      `json:"vod_score"`
		VodScoreAll      int         `json:"vod_score_all"`
		VodScoreNum      int         `json:"vod_score_num"`
		VodTime          string      `json:"vod_time"`
		VodTimeAdd       int         `json:"vod_time_add"`
		VodTimeHits      int         `json:"vod_time_hits"`
		VodTimeMake      int         `json:"vod_time_make"`
		VodTrysee        int         `json:"vod_trysee"`
		VodDoubanId      int         `json:"vod_douban_id"`
		VodDoubanScore   string      `json:"vod_douban_score"`
		VodReurl         string      `json:"vod_reurl"`
		VodRelVod        string      `json:"vod_rel_vod"`
		VodRelArt        string      `json:"vod_rel_art"`
		VodPwd           string      `json:"vod_pwd"`
		VodPwdUrl        string      `json:"vod_pwd_url"`
		VodPwdPlay       string      `json:"vod_pwd_play"`
		VodPwdPlayUrl    string      `json:"vod_pwd_play_url"`
		VodPwdDown       string      `json:"vod_pwd_down"`
		VodPwdDownUrl    string      `json:"vod_pwd_down_url"`
		VodContent       string      `json:"vod_content"`
		VodPlayFrom      string      `json:"vod_play_from"`
		VodPlayServer    string      `json:"vod_play_server"`
		VodPlayNote      string      `json:"vod_play_note"`
		VodPlayUrl       string      `json:"vod_play_url"`
		VodDownFrom      string      `json:"vod_down_from"`
		VodDownServer    string      `json:"vod_down_server"`
		VodDownNote      string      `json:"vod_down_note"`
		VodDownUrl       string      `json:"vod_down_url"`
		VodPlot          int         `json:"vod_plot"`
		VodPlotName      string      `json:"vod_plot_name"`
		VodPlotDetail    string      `json:"vod_plot_detail"`
		VodPicScreenshot interface{} `json:"vod_pic_screenshot"`
		TypeName         string      `json:"type_name"`
	} `json:"list"`
}
type Res struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
	Url  string `json:"url"`
}

func init() {
	rand.Seed(time.Now().Unix())
}
func clients(num int, d bool, randPage bool) {
	countPage := 0
	list, err := getList(1)
	if err != nil {
		logs.Info(err)
		return
	}
	countPage = list.Pagecount
	if randPage {
		i := rand.Int63n(int64(countPage)) + 1
		l, err := getList(int(i))
		if err == nil {
			Start(l, num, d)
		}
	}

}
func Start(list List, num int, d bool) {
	for _, s := range list.List {
		info, err := getInfo(s.VodId)
		if err != nil {
			logs.Info(err)
			return
		}
		split := strings.Split(info.List[0].VodPlayUrl, "$$$")
		if len(split) == 2 {
			i2 := strings.Split(split[1], "$")
			wp := workpool.New(num)
			for _, s2 := range i2 {
				if strings.Contains(s2, "http") {
					i3 := strings.Split(s2, "#")
					wp.Do(func() error {
						tr := &http.Transport{
							TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						}
						//http cookie接口
						cookieJar, _ := cookiejar.New(nil)
						c := &http.Client{
							Jar:       cookieJar,
							Transport: tr,
						}
						now := time.Now()
						get, err := c.Get("https://gpgo.site/get_new?url=" + i3[0])
						if err != nil {
							logs.Info(err)
							return nil
						}
						all, err := ioutil.ReadAll(get.Body)
						if err != nil {
							logs.Info(err)
							return nil
						}
						var res Res
						_ = json.Unmarshal(all, &res)
						if res.Err != "" {
							logs.Info(res.Err)
							return nil
						} else {
							logs.Info("获取url", i3[0], time.Now().Sub(now).String())
							if d {
								download(res.Url)
							}
						}
						return nil
					})
				}
			}
			wp.Wait()
		}
	}
}
func download(url string) {
	exec.Command("sh", "-c", "./m3u8D -u "+url+" -ht=\"apiv2\"").Run()
	//logs.Info("开始下载", url)
	////./m3u8D -u=https://m8t.vboku.com/20211007/T2QV9wNJ/hls/index.m3u8 -ht=apiv2
	//theURL, err := ParseM3U8AndCacheTheURL(url)
	//if err != nil {
	//	logs.Err(err, url)
	//} else {
	//	size := 0
	//	l := len(theURL)
	//	for i := range theURL {
	//		bytes, err2 := Download(theURL[l-i-1])
	//		if err2 != nil {
	//			logs.Err(err2)
	//		}
	//		size += len(bytes)
	//	}
	//	logs.Info("下载完成", fmt.Sprintf("%dMB", size/1024/1024))
	//}
}

//解析m3u8文件链接
func getM3U8UrlContent(m3u8Url string) (*url.URL, interface{}, m3u8s.ListType, error) {
	urlParse, err := url.Parse(m3u8Url)
	if err != nil {
		return nil, nil, 0, err
	}
	resp, err := Get(m3u8Url)
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
func ParseM3U8AndCacheTheURL(m3u8 string) ([]string, error) {
	urls := make([]string, 0)
	var err error
	urlP, list, listType, err := getM3U8UrlContent(m3u8)
	if err != nil {
		return nil, err
	}
	if listType == m3u8s.MASTER {
		playlist := list.(*m3u8s.MasterPlaylist)
		variants := playlist.Variants
		if len(variants) == 0 {
			return nil, errors.New("未获取到播放列表")
		}
		for i, variant := range variants {
			if !strings.Contains(variant.URI, "://") {
				variant.URI = HostAddPath(urlP, variant.URI)
			}
			urls = append(urls, variant.URI)
			theURL, err := parseMediaM3U8AndCacheTheURL(variant.URI, i)
			if err != nil {
				return nil, err
			}
			urls = append(urls, theURL...)
		}
		return urls, nil
	}
	theURL, err := parseMediaM3U8AndCacheTheURL(m3u8, 0)
	if err != nil {
		return nil, err
	}
	return append(urls, theURL...), nil
}
func HostAddPath(url *url.URL, path string) string {
	base := filepath.Dir(url.Path)
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	if base != "/" && !strings.HasPrefix(path, base) {
		path = base + path
	}
	if strings.Contains(url.Host, ":") {
		return fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, path)
	}
	return fmt.Sprintf("%s://%s:%s%s", url.Scheme, url.Host, url.Port(), path)
}
func parseMediaM3U8AndCacheTheURL(m3u8 string, index int) ([]string, error) {
	urls := make([]string, 0)
	urlP, list, listType, err := getM3U8UrlContent(m3u8)
	if err != nil {
		return nil, err
	}
	if listType != m3u8s.MEDIA {
		return nil, errors.New("类型错误 listType!=m3u8s.MEDIA")
	}
	mediaList := list.(*m3u8s.MediaPlaylist)
	segments := mediaList.Segments
	if len(segments) == 0 {
		return nil, errors.New("未获取到播放列表资源")
	}
	//存在加密
	if mediaList.Key != nil && mediaList.Key.URI != "" {
		keyUrl := mediaList.Key.URI
		if !strings.Contains(keyUrl, "://") {
			keyUrl = HostAddPath(urlP, keyUrl)
		}
		urls = append(urls, keyUrl)
	}
	//缓存ts url
	for _, segment := range segments {
		if segment == nil {
			continue
		}
		if !strings.Contains(segment.URI, "://") {
			segment.URI = HostAddPath(urlP, segment.URI)
		}
		urls = append(urls, segment.URI)
	}
	return urls, nil
}
