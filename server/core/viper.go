package core

import (
	"flag"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	_ "github.com/flipped-aurora/gin-vue-admin/server/packfile"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	buildTime      = "2023-01-01T00:00:00"
	currentVersion = "first commit"
)

func Viper(path ...string) *viper.Viper {
	var config string
	//info, _ := debug.ReadBuildInfo()
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		version := flag.Bool("v", false, "building info.")
		flag.Parse()
		if *version {
			//for _, setting := range info.Settings {
			//	fmt.Println(setting)
			//}
			fmt.Println("build time:", buildTime)
			fmt.Println("version:  ", currentVersion)
			return nil
		} else {
			if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
				if configEnv := os.Getenv(utils.ConfigEnv); configEnv == "" {
					config = utils.ConfigFile
					fmt.Printf("您正在使用config的默认值,config的路径为%v\n", utils.ConfigFile)
				} else {
					config = configEnv
					fmt.Printf("您正在使用GVA_CONFIG环境变量,config的路径为%v\n", config)
				}
			} else {
				fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
			}
		}
	} else {
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	// root 适配性
	// 根据root位置去找到对应迁移位置,保证root路径有效
	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(global.GVA_CONFIG.JWT.ExpiresTime)),
	)
	return v
}
