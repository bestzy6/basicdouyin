package server

import (
	"basictiktok/api"
	"basictiktok/server/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改（此处中间件是gin中间件）
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// 路由
	v1 := r.Group("/douyin")
	{
		//基础接口
		v1.GET("/feed", api.Ping)            //视频流接口
		v1.POST("/user/register", api.Ping)  //用户注册
		v1.POST("/user/login", api.Ping)     //用户登录
		v1.GET("/user", api.Ping)            //用户信息
		v1.POST("/publish/action", api.Ping) //投稿接口
		v1.GET("/publish/list", api.Ping)    //发布列表

		//拓展接口1
		v1.POST("/favorite/action", api.Ping) //赞操作
		v1.GET("/favorite/list", api.Ping)    //点赞列表
		v1.POST("/comment/action", api.Ping)  //评论操作
		v1.GET("/comment/list", api.Ping)     //评论列表

		//拓展接口2
		v1.POST("/relation/action", api.Ping)       //关注操作
		v1.GET("/relation/follow/list", api.Ping)   //关注列表
		v1.GET("/relation/follower/list", api.Ping) //粉丝列表

		//以下是脚手架自带,供参考，之后会删除
		v1.POST("ping", api.Ping)
	}
	return r
}
