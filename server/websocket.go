package server

import (
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"ws/logic"
)

// WebSocketHandleFunc 服务处理
func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 从客户端接收 Websocket 握手，并将连接升级到 Websocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项 （AcceptOptions来设置）

	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 1.新用户进来，构建该用户的实例
	nickname := req.FormValue("nickname")
	token := req.FormValue("token")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal:", nickname)
		log.Printf("%T", logic.NewErrorMessage("昵称已经存在"))
		err = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称"))
		if err != nil {
			log.Println(err)
		}
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal！")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {

		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("昵称已经存在"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists！")
		return
	}

	user := logic.NewUser(conn, token, nickname, req.RemoteAddr)
	log.Println("初始化用户：", user)
	// 2.开启给用户发送消息的goroutine
	go user.SendMessage(req.Context())

	// 3.给新用户发送欢迎消息
	user.MessageChannel <- logic.NewWelcomeMessage(user)

	// 向所有用户告知新用户的到来
	msg := logic.NewEnterMessage(user)
	logic.Broadcaster.BroadCaster(msg)
	log.Println(msg)

	// 4.将用户加入广播器列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	// 5.接收用户发送过来的信息
	err = user.ReceiveMessage(req.Context())

	// 6.用户离开
	logic.Broadcaster.UserLeaving(user)
	// 向所有用户告知用户离开
	msg = logic.NewLeaveMessage(user)
	logic.Broadcaster.BroadCaster(msg)
	log.Println(msg)

	// 根据读取时的错误执行不同的 Close

	if err != nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read room client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
