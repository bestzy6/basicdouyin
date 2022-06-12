package server

import (
	"basictiktok/api"
	"basictiktok/server/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/static/video", "./static/video")
	r.Static("/static/img", "./static/img")
	// 中间件
	r.Use(middleware.Cors())

	// 路由
	v1 := r.Group("/douyin")
	{
		//token鉴定（需要鉴权的接口把前面改为needAuth）
		needAuth := v1.Group("")
		needAuth.Use(middleware.AuthToken())

		//基础接口
		v1.GET("/feed", api.VideoList)               //视频流接口
		v1.POST("/user/register/", api.Register)     //用户注册
		v1.POST("/user/login", api.Login)            //用户登录
		v1.GET("/user", api.QueryUserInfo)           //用户信息
		v1.POST("/publish/action", api.PublishVideo) //投稿接口
		v1.GET("/publish/list/", api.ListVideos)     //发布列表

		//拓展接口1
		v1.POST("/favorite/action", api.Ping) //赞操作
		v1.GET("/favorite/list", api.Ping)    //点赞列表
		v1.POST("/comment/action", api.Ping)  //评论操作
		v1.GET("/comment/list", api.Ping)     //评论列表

		//拓展接口2
		needAuth.POST("/relation/action", api.FollowAction)       //关注操作
		needAuth.GET("/relation/follow/list", api.GetFollowers)   //关注列表
		needAuth.GET("/relation/follower/list", api.GetFollowees) //粉丝列表
	}
	return r
}
