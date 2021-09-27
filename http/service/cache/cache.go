package cache

import (
	"errors"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"io/ioutil"
	"os"
	downloadServer "p00q.cn/video_cdn/http/service/download"
	"p00q.cn/video_cdn/http/utils"
	"path/filepath"
	"time"
)

var cacheDir = "./cache"
var cacheKeyMap = make(map[string]string)

func init() {
	if !utils.IsAbsPath(cacheDir) {
		abs, err := utils.Abs(cacheDir)
		if err != nil {
			logs.Err("缓存目录路径在转换绝对路径时遇到问题", err)
		} else {
			cacheDir = abs
		}
	}
	err := utils.IsDirExistCreateIt(cacheDir)
	if err != nil {
		logs.Err(err)
	}
	err = utils.IsDirExistCreateIt(filepath.Join(cacheDir, "tmp"))
	if err != nil {
		logs.Err(err)
	}
}
func CacheKey(key string, keyVal string) {
	cacheKeyMap[key] = keyVal
}
func GetCacheKey(key string) string {
	if v, has := cacheKeyMap[key]; has {
		return v
	}
	return ""
}
func CacheFormUrl(url string) ([]byte, error) {
	now := time.Now()
	data, err := downloadServer.Download(url)
	seconds1 := time.Now().Sub(now).Seconds()
	if err != nil {
		return nil, err
	}
	go Cache(url, data)
	seconds2 := time.Now().Sub(now).Seconds()
	logs.Info(fmt.Sprintf("%.2f %.2f", seconds1, seconds2))
	return data, nil
}
func Cache(key string, data []byte) error {
	md5 := utils.MD5(key)
	cachePath := getPath(md5)
	err := utils.IsDirExistCreateIt(filepath.Dir(cachePath))
	if err != nil {
		return err
	}
	exists := utils.Exists(cachePath)
	if exists {
		return errors.New("已存在缓存")
	}
	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	defer file.Close()
	if err != nil {
		return err
	}
	return err
}
func HasKey(key string) bool {
	md5 := utils.MD5(key)
	cachePath := getPath(md5)
	return utils.Exists(cachePath)
}
func GetCache(key string) ([]byte, error) {
	md5 := utils.MD5(key)
	cachePath := getPath(md5)
	return ioutil.ReadFile(cachePath)
}
func getDir(md5 string) string {
	if len(md5) != 32 {
		return ""
	}
	return filepath.Join(cacheDir, md5[:3])
}
func getPath(md5 string) string {
	if len(md5) != 32 {
		return ""
	}
	return filepath.Join(cacheDir, md5[:3], md5)
}
