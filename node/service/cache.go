package service

import (
	"errors"
	logs "github.com/danbai225/go-logs"
	"github.com/shirou/gopsutil/v3/disk"
	"io/fs"
	"io/ioutil"
	"os"
	"p00q.cn/video_cdn/comm/utils"
	"p00q.cn/video_cdn/node/config"
	"path/filepath"
	"strings"
	"time"
)

func CacheFormUrl(url string) ([]byte, error) {
	data, err := Download(GetUrl(url))
	if err != nil {
		return nil, err
	}
	go Cache(url, data)
	return data, nil
}
func Cache(key string, data []byte) error {
	if len(data) < 1024*20 && strings.HasSuffix(key, "ts") {
		return errors.New("size <20k")
	}
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
	var usage *disk.UsageStat
	var err error
	if clearFlg {
		return
	} else {
		clearFlg = true
		defer func() {
			clearFlg = false
			logs.Info("清理完成", usage.UsedPercent, usage.Used/1024/1024, "MB")
		}()
	}
	usage, err = disk.Usage(config.GlobalConfig.CacheDir)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info("开始清理", usage.UsedPercent, usage.Used/1024/1024, "MB")
	diffTime := int64(3600 * 24 * 7)
	now := time.Now().Unix()
	for usage.UsedPercent > 80 {
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
		usage, err = disk.Usage(config.GlobalConfig.CacheDir)
		if err != nil {
			logs.Err(err)
			return
		}
		diffTime = diffTime / 2
	}
}
