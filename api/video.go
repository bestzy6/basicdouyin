package api

import (
	"basictiktok/serializer"
	"basictiktok/server/middleware"
	"basictiktok/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

//

func PublishVideo(c *gin.Context) {
	var req serializer.ActionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	userid, _ := c.Get("userid")
	resp := service.ActionService(&req, userid.(int), c.Request.Host)
	c.JSON(http.StatusOK, resp)
}

// VideoList 视频流接口
func VideoList(c *gin.Context) {
	var req serializer.FeedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	userid := 0
	token, exist := c.GetQuery("token")
	if exist {
		claims, err := middleware.ParseToken(token)
		if err == nil {
			userid = claims.UserID
		}
	}
	resp := service.FindVideoBeforeTimeService(&req, userid)
	c.JSON(http.StatusOK, resp)
}

// ListVideos 用户视频发布列表
func ListVideos(c *gin.Context) {
	var req serializer.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	id, _ := c.Get("userid")
	userid := id.(int)
	resp := service.ListVideosService(&req, userid)
	c.JSON(http.StatusOK, resp)
}
