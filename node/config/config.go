package config

import (
	"encoding/json"
	logs "github.com/danbai225/go-logs"
	"io/ioutil"
	"p00q.cn/video_cdn/node/utils"
	"path/filepath"
)

var GlobalConfig Config
var Path = "./config.json"

type Config struct {
	ServerAddress string
	Token         string
	Port          int
	CacheDir      string
	Host          string
}

func LoadConfig() error {
	config := Config{}
	data, err := ioutil.ReadFile(Path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	if config.CacheDir == "" {
		config.CacheDir = "./cache"
	}
	GlobalConfig = config
	initFunc()
	return nil
}
func initFunc() {
	if !utils.IsAbsPath(GlobalConfig.CacheDir) {
		abs, err := utils.Abs(GlobalConfig.CacheDir)
		if err != nil {
			logs.Err("缓存目录路径在转换绝对路径时遇到问题", err)
		} else {
			GlobalConfig.CacheDir = abs
		}
	}
	err := utils.IsDirExistCreateIt(GlobalConfig.CacheDir)
	if err != nil {
		logs.Err(err)
	}
	err = utils.IsDirExistCreateIt(filepath.Join(GlobalConfig.CacheDir, "tmp"))
	if err != nil {
		logs.Err(err)
	}
}
