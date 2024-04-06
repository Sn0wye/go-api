package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func GetConfig() *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	// TODO: Uncomment this code if you want to use flag to set the config path
	// if envConf == "" {
	// 	fmt.Println("No config path provided, using default config path")
	// 	flag.StringVar(&envConf, "conf", "config/local.yml", "config path, eg: -conf config/local.yml")
	// 	flag.Parse()
	// }
	if envConf == "" {
		envConf = "../config/local.yml"
	}

	basepath := getConfigPath() + "/config/local.yml"

	return getConfig(basepath)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

func getConfigPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	index := strings.LastIndex(basepath, "/pkg")
	if index != -1 {
		basepath = basepath[:index]
	}

	return basepath
}
