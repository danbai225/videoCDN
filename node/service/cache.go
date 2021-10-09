package service

import (
	"errors"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/shirou/gopsutil/v3/disk"
	"io/fs"
	"io/ioutil"
	"os"
	"p00q.cn/video_cdn/comm/utils"
	"p00q.cn/video_cdn/node/config"
	"path/filepath"
	"time"
)

func CacheFormUrl(url string) ([]byte, error) {
	now := time.Now()
	data, err := Download(GetUrl(url))
	logs.Info("cache", url, fmt.Sprintf("%0.2f", time.Now().Sub(now).Seconds()))
	if err != nil {
		return nil, err
	}
	go Cache(url, data)
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
	return filepath.Join(config.GlobalConfig.CacheDir, md5[:3])
}
func getPath(md5 string) string {
	if len(md5) != 32 {
		return ""
	}
	return filepath.Join(config.GlobalConfig.CacheDir, md5[:3], md5)
}

var clearFlg bool

//清理磁盘缓存
func clear() {
	if clearFlg {
		return
	} else {
		clearFlg = true
		defer func() {
			clearFlg = false
			logs.Info("清理完成")
		}()
	}
	logs.Info("开始清理")
	diffTime := int64(3600 * 24 * 7)
	UsedPercent := 81
	now := time.Now().Unix()
	for UsedPercent > 80 {
		_ = filepath.WalkDir(config.GlobalConfig.CacheDir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				info, err := d.Info()
				if err == nil {
					fileTime := info.ModTime().Unix()
					if (now - fileTime) > diffTime {
						os.Remove(path)
					}
				}
			}
			return nil
		})
		usage, err := disk.Usage(config.GlobalConfig.CacheDir)
		if err == nil {
			UsedPercent = int(usage.UsedPercent)
		} else {
			return
		}
		diffTime = diffTime / 2
	}
}
