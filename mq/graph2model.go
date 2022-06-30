package mq

import "basictiktok/model"

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
	var msg *UserMessage
	for {
		msg = <-ToModelUserMQ
		user, toUser, video := msg.User, msg.ToUser, msg.ToVideo
		switch msg.OpNum {
		case Follow:
			user.Follow(toUser)
		case UnFollow:
			user.UnFollow(toUser)
		case Favorite:
			videoDao := model.NewVideoDaoInstance()
			videoDao.AddFavorite(video.ID)
		case UnFavorite:
			videoDao := model.NewVideoDaoInstance()
			videoDao.DeleteFavorite(video.ID)
		}
	}
}
