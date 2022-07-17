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
	go listenFollowMQ(mq.FollowConsumerMsg)
	go listenFavoriteMQ(mq.FavoriteConsumerMsg)
}

func listenFollowMQ(msg <-chan string) {
	userDaoInstance := dao.NewUserDaoInstance()
	for {
		str := <-msg
		split := strings.Split(str, "_")
		userID, _ := strconv.Atoi(split[0])
		targetUserId, _ := strconv.Atoi(split[1])
		actionType, _ := strconv.Atoi(split[2])
		var err error
		if actionType == 1 {
			err = userDaoInstance.Follow(userID, targetUserId)
		} else {
			err = userDaoInstance.UnFollow(userID, targetUserId)
		}
		if err != nil {
			util.Log().Error("MQ err:", err)
		}
	}
}

func listenFavoriteMQ(msg <-chan string) {
	videoDaoInstance := dao.NewVideoDaoInstance()
	for {
		str := <-msg
		split := strings.Split(str, "_")
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
	}
}
