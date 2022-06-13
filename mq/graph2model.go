package mq

import "basictiktok/model"

//通道的缓存最大值
const maxMessageNum = 100

type OperaNum uint8

const (
	Follow   OperaNum = iota //关注
	UnFollow                 //取消关注
)

type G2mMessage struct {
	User   *model.User
	ToUser *model.User
	Num    OperaNum
}

var (
	ToModelUserMQ chan *G2mMessage
)

func InitMQ() {
	ToModelUserMQ = make(chan *G2mMessage, maxMessageNum)
	go listenToModelUserMQ()
}

func listenToModelUserMQ() {
	var msg *G2mMessage
	for {
		msg = <-ToModelUserMQ
		user, toUser := msg.User, msg.ToUser
		switch msg.Num {
		case Follow:
			user.Follow(toUser)
		case UnFollow:
			user.UnFollow(toUser)
		}
	}
}
