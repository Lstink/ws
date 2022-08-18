package main

import (
	"log"
	"net/http"
	"ws/global"
	"ws/server"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

Go 语言编程之旅 —— 一起用 Go 做项目：ChatRoom，start on：%s
`
)

func init() {
	global.Init()
}

func main() {

	log.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
