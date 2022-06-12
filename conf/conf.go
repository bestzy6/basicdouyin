package conf

import (
	"basictiktok/cache"
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/util"
	"os"

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
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
	graphdb.Neo4j()
}
