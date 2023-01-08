package logic

import (
	"time"

	"github.com/spf13/cast"
)

const (
	MsgTypeNormal    = iota // 普通 用户訊息
	MsgTypeWelcome          // 當前用户歡迎訊息
	MsgTypeUserEnter        // 用户進入
	MsgTypeUserLeave        // 用户退出
	MsgTypeError            // 错误訊息
)

// 给用户發送的消息
type Message struct {
	// 哪個用户發送的訊息
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	// 訊息 @ 了谁
	Ats []string `json:"ats"`

	// 用户列表不通過WebSocket 下發
	// Users []*User `json:"users"`
}

func NewMessage(user *User, content, clientTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return message
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: user.NickName + " 歡迎加入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.NickName + " 離開了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}