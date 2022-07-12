package mq

import (
	"basictiktok/dao"
	"basictiktok/model"
)

//通道的缓存最大值
const maxMessageNum = 100

type OperaNum uint8

const (
	Follow     OperaNum = iota //关注
	UnFollow                   //取消关注
	Favorite                   //点赞
	UnFavorite                 //取消点赞
)

type UserMessage struct {
	User    *model.User
	ToUser  *model.User
	ToVideo *model.Video
	OpNum   OperaNum
}

var (
	ToModelUserMQ chan *UserMessage
)

func InitMQ() {
	ToModelUserMQ = make(chan *UserMessage, maxMessageNum)
	go listenToModelUserMQ()
}

func listenToModelUserMQ() {
	var (
		msg      *UserMessage
		userDao  = dao.NewUserDaoInstance()
		videoDao = dao.NewVideoDaoInstance()
	)
	for {
		msg = <-ToModelUserMQ
		user, toUser, video := msg.User, msg.ToUser, msg.ToVideo
		switch msg.OpNum {
		case Follow:
			userDao.Follow(user.ID, toUser.ID)
		case UnFollow:
			userDao.UnFollow(user.ID, toUser.ID)
		case Favorite:
			videoDao.AddFavorite(video.ID)
		case UnFavorite:
			videoDao.DeleteFavorite(video.ID)
		}
	}
}
