package server

import (
	"net/http"
	"os"
	"path/filepath"
	"ws/logic"
)

func RegisterHandle() {
	inferRootDir()

	// 广播消息处理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

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

	rootDir = infer(cwd)
}

func exists(s string) bool {
	_, err := os.Stat(s)
	return os.IsNotExist(err)
}

var rootDir string
