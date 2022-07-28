package conf

import (
	"basictiktok/cache"
	"basictiktok/dao"
	"basictiktok/graphdb"
	"basictiktok/mq"
	"basictiktok/util"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	err := godotenv.Load()
	if err != nil {
		util.Log().Panic("读取.env文件失败！")
		return
	}

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 连接数据库
	dao.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
	graphdb.Neo4j()

	//消息队列
	mq.InitKafka()

	//监控消息队列
	go listenFollowMQ(mq.FollowConsumerMsg, mq.FollowNotifyMsg)
	go listenFavoriteMQ(mq.FavoriteConsumerMsg, mq.FavoriteNotifyMsg)
}

func listenFollowMQ(msg <-chan string, notify chan<- struct{}) {
	userDaoInstance := dao.NewUserDaoInstance()
	offsetId := 0
	for {
		str := <-msg
		split := strings.Split(str, "_")
		snowId, _ := strconv.Atoi(split[0])
		if snowId > offsetId {
			userID, _ := strconv.Atoi(split[1])
			targetUserId, _ := strconv.Atoi(split[2])
			actionType, _ := strconv.Atoi(split[3])
			var err error
			if actionType == 1 {
				err = userDaoInstance.Follow(userID, targetUserId)
			} else {
				err = userDaoInstance.UnFollow(userID, targetUserId)
			}
			if err != nil {
				util.Log().Error("MQ err:", err)
			}
			offsetId = snowId
		}
		//表示已经执行完成
		notify <- struct{}{}
	}
}

func listenFavoriteMQ(msg <-chan string, notify chan<- struct{}) {
	videoDaoInstance := dao.NewVideoDaoInstance()
	offsetId := 0
	for {
		str := <-msg
		split := strings.Split(str, "_")
		currentId, _ := strconv.Atoi(split[0])
		if currentId > offsetId {
			videoId, _ := strconv.Atoi(split[1])
			actionType, _ := strconv.Atoi(split[2])
			var err error
			if actionType == 1 {
				err = videoDaoInstance.AddFavorite(int64(videoId))
			} else {
				err = videoDaoInstance.DeleteFavorite(int64(videoId))
			}
			if err != nil {
				util.Log().Error("MQ err:", err)
			}
			offsetId = currentId
		}
		//表示已经执行完成
		notify <- struct{}{}
	}
}
