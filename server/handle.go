package server

import (
	"net/http"
	"ws/logic"
)

func RegisterHandle() {

	// 广播消息处理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
}
