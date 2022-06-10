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
	user     *model.User
	operaNum OperaNum
}

var (
	toModelUserMQ chan *G2mMessage
)

func InitMQ() {
	toModelUserMQ = make(chan *G2mMessage, maxMessageNum)
	go listenToModelUserMQ()
}

func listenToModelUserMQ() {
	var message *G2mMessage
	for {
		message = <-toModelUserMQ
		user := message.user
		switch message.operaNum {
		case DecreFollower:

			return
		case IncreFollower:
			return
		case DecreFollowee:
			return
		case IncreFollowee:
			return
		}
	}
}
