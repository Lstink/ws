package logic

import (
	"context"
	"errors"
	"io"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nick_name"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"_"`

	conn *websocket.Conn
}

// 系统用户，代表是系统发送的消息
var System = &User{}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)

	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// 判断连接是否关闭了，正常关闭，不认为是错误
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				// 如果是关闭
				return nil
			} else if errors.Is(err, io.EOF) {
				// 如果是发送结束
				return nil
			}

			return err
		}

		// 内容发送到聊天室
		sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		// 过滤消息内容
		sendMsg.Content = FilterSensitive(sendMsg.Content)

	}

}

func NewUser(conn *websocket.Conn, nickname, addr string) *User {
	return &User{conn: conn, NickName: nickname, Addr: addr}
}
