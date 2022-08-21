package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

var globalUID uint32 = 0

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nick_name"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	Token          string        `json:"token"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn

	isNew bool
}

// System 系统用户，代表是系统发送的消息
var System = &User{}

// SendMessage 给用户发送消息
func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		log.Println("给用户发送消息：", msg, u)
		// 给用户发送消息
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
		log.Println("将要发送的消息", sendMsg)

		// 过滤消息内容
		sendMsg.Content = FilterSensitive(sendMsg.Content)
		// 解析 content，看看 @ 谁了
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)

		Broadcaster.BroadCaster(sendMsg)

	}

}

func NewUser(conn *websocket.Conn, token, nickname, addr string) *User {
	user := &User{
		conn:           conn,
		NickName:       nickname,
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		EnterAt:        time.Now(),
		Token:          token,
	}

	if user.Token != "" && user.Token != "undefined" {
		log.Println("token是什么", user.Token)
		uid, err := parseTokenAndValidate(token, nickname)
		if err == nil {
			user.UID = uid
		}
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
		user.Token = genToken(user.UID, user.NickName)
		user.isNew = true
	}

	return user
}

func genToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)

	messageMAC := macSha256([]byte(message), []byte(secret))

	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

func parseTokenAndValidate(token, nickname string) (int, error) {
	pos := strings.LastIndex(token, "uid")
	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)

	ok := validateMAC([]byte(message), messageMAC, []byte(secret))
	if ok {
		return uid, nil
	}

	return 0, errors.New("token is illegal")
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func validateMAC(message, messageMAC, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
