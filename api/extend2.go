package api

import (
	"basictiktok/serializer"
	"basictiktok/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FollowAction 关注或取消关注
func FollowAction(c *gin.Context) {
	var req serializer.FollowRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.FollowResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	reqUserId, _ := c.Get("userid")
	req.ReqUserId = reqUserId.(int)
	var resp *serializer.FollowResponse
	if req.ActionType == 1 {
		resp = service.FollowService(&req)
		c.JSON(http.StatusOK, resp)
		return
	}
	if req.ActionType == 2 {
		resp = service.UnFollowService(&req)
		c.JSON(http.StatusOK, resp)
		return
	}
	c.JSON(http.StatusOK, serializer.FollowResponse{
		StatusCode: serializer.ParamInvalid,
		StatusMsg:  "请求参数错误",
	})
}

// GetFollowers 获取关注列表
func GetFollowers(c *gin.Context) {
	var req serializer.FollowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.FollowResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	reqUserId, _ := c.Get("userid")
	req.ReqUserId = reqUserId.(int)
	resp := service.FollowersService(&req)
	c.JSON(http.StatusOK, resp)
}

// GetFollowees 获取粉丝列表
func GetFollowees(c *gin.Context) {
	var req serializer.FolloweesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.FollowResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	reqUserId, _ := c.Get("userid")
	req.ReqUserId = reqUserId.(int)
	resp := service.FolloweesService(&req)
	c.JSON(http.StatusOK, resp)
}
