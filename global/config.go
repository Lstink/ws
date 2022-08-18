package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords []string

	MessageQueueLen = 1024
)

func initConfig() {
	// 设置读取配置文件的名称
	viper.SetConfigName("chatroom")
	viper.SetConfigType("yaml")
	// 设置配置文件所在目录
	viper.AddConfigPath(RootDir + "/config")
	// 如果读取配置文件错误，则 panic
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 读取指定key内容为字符切片
	SensitiveWords = viper.GetStringSlice("sensitive")
	// 读取 int 类型变量
	MessageQueueLen = viper.GetInt("message-queue")

	// 监控文件
	viper.WatchConfig()
	// 如果文件内容发送变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 读取文件
		viper.ReadInConfig()
		// 重新赋值
		SensitiveWords = viper.GetStringSlice("sensitive")
	})

}
