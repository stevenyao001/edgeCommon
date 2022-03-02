package config

import (
	"github.com/spf13/viper"
)


func InitConf(filePath string,conf interface{}) {
	//设置配置文件类型
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		panic("init config fail :"+err.Error())
	}

	if err := viper.Unmarshal(conf); err != nil {
		panic("resolution config fail :"+err.Error())
	}
}
