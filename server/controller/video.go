package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"io/ioutil"
	"net/http"
	"strings"
)

func infoH(r *ghttp.Request) {
	id := r.GetQueryInt("id", 1)
	l, err := getInfo(id)
	if err != nil {
		r.Response.WriteJson(g.Map{"err": err.Error(), "code": 1})

	} else {
		ji(&l)
		r.Response.WriteJson(g.Map{"err": "", "code": 0, "info": l})
	}
}
func ji(info *Info) {
	if len(info.List) < 1 {
		return
	}
	info.List[0].Ji = make([]string, 0)
	split := strings.Split(info.List[0].VodPlayUrl, "$$$")
	if len(split) == 2 {
		i2 := strings.Split(split[1], "$")
		for _, s2 := range i2 {
			if strings.Contains(s2, "http") {
				info.List[0].Ji = append(info.List[0].Ji, strings.Split(s2, "#")[0])
			}
		}
	}
}
func listH(r *ghttp.Request) {
	page := r.GetQueryInt("page", 1)
	key := r.GetQueryString("key", "")
	l, err := getList(page, key)
	if err != nil {
		r.Response.WriteJson(g.Map{"err": err.Error(), "code": 1})

	} else {
		r.Response.WriteJson(g.Map{"err": "", "code": 0, "list": l})
	}
}

func getList(pg int, key string) (List, error) {
	get, err := http.Get(fmt.Sprintf("https://api.apibdzy.com/api.php/provide/vod/?ac=list&pg=%d&wd=%s", pg, key))
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
		Ji               []string    `json:"ji"`
	} `json:"list"`
}
