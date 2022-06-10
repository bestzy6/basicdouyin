package mq

import "basictiktok/model"

//通道的缓存最大值
const maxMessageNum = 100

type OperaNum uint8

const (
	DecreFollower OperaNum = iota //关注人-1
	IncreFollower                 //关注人+1
	DecreFollowee                 //粉丝-1
	IncreFollowee                 //粉丝+1
)

type G2mMessage struct {
	User *model.User
	Num  OperaNum
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
		user := msg.User
		switch msg.Num {
		case DecreFollower:
			user.DecreFollow()
		case IncreFollower:
			user.IncreFollow()
		case DecreFollowee:
			user.DecreFollowee()
		case IncreFollowee:
			user.IncreFollowee()
		}
	}
}
