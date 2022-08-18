package global

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string

var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

// 推测项目根目录
func inferRootDir() {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	var infer func(d string) string

	infer = func(d string) string {

		// 这里要确保项目目录下存在 template 目录
		if exists(d + "/template") {
			return d
		}
		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
	fmt.Println("推测的目录", RootDir)
}

// 判断目录是否存在
func exists(s string) bool {
	_, err := os.Stat(s)
	return !os.IsNotExist(err)
}
