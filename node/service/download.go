package service

import (
	"errors"
	"github.com/levigross/grequests"
	"time"
)

var (
	requestOptions = &grequests.RequestOptions{
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
		RequestTimeout: time.Second * 10,
		Headers: map[string]string{
			"Connection":      "keep-alive",
			"Accept":          "*/*",
			"Accept-Encoding": "*",
			"Accept-Language": "zh-Hans;q=1",
			"Host":            "gpgo.site",
			"Referer":         "https://gpgo.site/",
		},
		InsecureSkipVerify: true,
	}
)

func Download(url string) ([]byte, error) {
	response, err := Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		//被限制下载尝试重定向源url给客户端
		if response.StatusCode >= 500 {
			return []byte(url), errors.New("redirect")
		}
		return nil, err
	}
	return response.Bytes(), err
}
func Get(url string) (*grequests.Response, error) {
	//logs.Info(url)
	response, err := grequests.Get(url, requestOptions)
	if err != nil {
		return nil, err
	}
	return response, nil
}
