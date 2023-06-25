package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Conf = new(Config)

func init() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("get pwd path failed:%v\n", err)
		return
	}
	filepath := path + CONFIG_DIR
	c := viper.New()
	c.SetConfigFile(CONFIG_FILE)
	c.AddConfigPath(filepath)         //设置读取的文件路径
	c.SetConfigName(CONFIG_FILE_NAME) //设置读取的文件名
	c.SetConfigType(CONFIG_TYPE)      //chaos设置文件的类型

	err = c.ReadInConfig() // 搜索并读取配置文件
	if err != nil {        // 处理错误
		fmt.Printf("read config file failed:%v\n", err)
		return
	}
	err = c.Unmarshal(&Conf) //将配置文件绑定到config上
	if err != nil {
		fmt.Printf("unmarshal config  failed:%v\n", err)
		return
	}

	fmt.Printf("init config success\n")
}
