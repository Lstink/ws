package logic

import (
	"context"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type Message struct {
	// 哪个用户发送的消息
	User    *User            `json:"user"`
	Type    int              `json:"type"`
	Content string           `json:"content"`
	MsgTime time.Time        `json:"msg_time"`
	Users   map[string]*User `json:"users"`
}

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nick_name"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"_"`

	conn *websocket.Conn
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func NewUser(conn *websocket.Conn, nickname, addr string) *User {
	return &User{conn: conn, NickName: nickname, Addr: addr}
}
