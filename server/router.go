package server

import (
	"basictiktok/api"
	"basictiktok/server/middleware"
	"basictiktok/util"
	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	//存放截图和地址的位置
	r.Static("/static/img", util.IMG)
	r.Static("/static/video", util.VEDIO)
	// 中间件
	r.Use(middleware.Cors())

	// 路由
	v1 := r.Group("/douyin")
	{
		//token鉴定（需要鉴权的接口把前面改为needAuth）
		needAuth := v1.Group("")
		needAuth.Use(middleware.AuthToken())

		//基础接口
		v1.GET("/feed", api.VideoList)                      //视频流接口
		v1.POST("/user/register/", api.Register)            //用户注册
		v1.POST("/user/login/", api.Login)                  //用户登录
		needAuth.GET("/user/", api.QueryUserInfo)           //用户信息
		needAuth.POST("/publish/action/", api.PublishVideo) //投稿接口
		needAuth.GET("/publish/list/", api.ListVideos)      //发布列表

		//拓展接口1
		needAuth.POST("/favorite/action/", api.FavoritePost) //赞操作
		needAuth.GET("/favorite/list/", api.FavoriteList)    //点赞列表
		needAuth.POST("/comment/action/", api.CommentPost)   //评论操作
		needAuth.GET("/comment/list/", api.CommentList)      //评论列表

		//拓展接口2
		needAuth.POST("/relation/action/", api.FollowAction)       //关注操作
		needAuth.GET("/relation/follow/list/", api.GetFollowers)   //关注列表
		needAuth.GET("/relation/follower/list/", api.GetFollowees) //粉丝列表
	}
	return r
}
