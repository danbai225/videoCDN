package main

import (
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
		},
	}
)

func Download(url string) ([]byte, error) {
	//./m3u8D -u http://107.189.14.25:8080/video/test/index.m3u8 -ht="apiv2"
	response, err := Get(url)
	if err != nil {
		return nil, err
	}
	return response.Bytes(), err
}
func Get(url string) (*grequests.Response, error) {
	response, err := grequests.Get(url, requestOptions)
	if err != nil {
		return nil, err
	}
	return response, nil
}
