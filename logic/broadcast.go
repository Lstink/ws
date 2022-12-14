package logic

import (
	"log"
	"ws/global"
)

// broadcaster 广播器
type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User

	//所有channel 统一管理，避免外部混用

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断用户是否可以进入聊天室
	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

// Broadcaster 单例模式，全局只使用一个对象
var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringChannel:       make(chan *User),
	leavingChannel:        make(chan *User),
	messageChannel:        make(chan *Message, global.MessageQueueLen),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
	requestUsersChannel:   make(chan struct{}),
	usersChannel:          make(chan []*User),
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) UserEntering(user *User) {
	log.Println("dayin", user)
	b.enteringChannel <- user
}

func (b *broadcaster) Start() {
	log.Println("启动消息管理器：")
	// 循环这个方法，让这个 goroutine 一直运行着
	for {
		select {
		case user := <-b.enteringChannel:
			// 新用户进入
			b.users[user.NickName] = user
		case user := <-b.leavingChannel:
			log.Println("检测到用户离开：", user)
			// 用户离开
			delete(b.users, user.NickName)
			// 避免 goroutine 泄漏
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			// 给所有在线用户发送消息
			for _, user := range b.users {
				// 如果是自己，就不发了
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}

	}
}

func (b *broadcaster) UserLeaving(user *User) {
	b.leavingChannel <- user
}

func (b *broadcaster) BroadCaster(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcaster queue is full!")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
