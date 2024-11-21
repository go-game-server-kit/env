package env_config_loader

import (
	"fmt"
	"github.com/go-game-server-kit/utils"
	"github.com/go-game-server-kit/xoe"
	"github.com/gookit/ini/v2/dotenv"
	"log"
	"os"
	"path"
)

type AppInfo struct {
	Name     string
	Env      string
	Group    string
	GroupTag string
	Version  string
	WorkDir  string
}

var (
	logger = log.New(os.Stderr, "[env] ", log.LstdFlags|log.Lshortfile)
)

var (
	info     *AppInfo
	infoOnce = utils.NewOnce(loadAppInfo)
	envOnce  = utils.NewOnce(loadEnv)
)

func loadAppInfo() *AppInfo {
	pwd := xoe.Poe1(os.Getwd())
	logger.Println("pwd:", pwd)

	defName := path.Base(pwd)
	info = &AppInfo{
		Name:     dotenv.Get("APP_NAME", defName),
		Env:      dotenv.Get("APP_ENV", "local"),
		Group:    dotenv.Get("APP_GROUP", ""),
		GroupTag: dotenv.Get("APP_GROUP_TAG", ""),
		Version:  dotenv.Get("APP_VERSION", ""),
		WorkDir:  pwd,
	}
	logger.Println("AppInfo:", utils.JsonPretty(info))
	return info
}

func loadEnv() struct{} {
	info := GetAppInfo()

	dir := info.WorkDir
	dirs := []string{dir}
	for {
		if dir == "/" {
			break
		}
		dir = path.Dir(dir)
		dirs = append(dirs, dir)
	}

	files := []string{".env"}
	if info.Env != "" {
		files = append([]string{fmt.Sprintf(".env.%s", info.Env)}, files...)
		if info.Group != "" {
			files = append([]string{fmt.Sprintf(".env.%s.%s", info.Env, info.Group)}, files...)
			if info.GroupTag != "" {
				files = append([]string{fmt.Sprintf(".env.%s.%s%s", info.Env, info.Group, info.GroupTag)}, files...)
			}
		}
	}

	for _, dir := range dirs {
		for _, file := range files {
			envFile := path.Join(dir, file)
			if _, err := os.Stat(envFile); err == nil {
				logger.Println("env load file:", envFile)
				xoe.Poe(dotenv.LoadFiles(envFile))
				goto end
			}
		}
	}
end:
	// env 会重写 env，所以需要重新加载
	logger.Println("reload app info")
	loadAppInfo()
	return struct{}{}
}

func LoadEnv() {
	envOnce.Do()
}

func GetAppInfo() *AppInfo {
	return infoOnce.Do()
}
