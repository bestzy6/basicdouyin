package api

import (
	"basictiktok/serializer"
	"basictiktok/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FavoritePost(c *gin.Context) {
	var req serializer.LikesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.CommentResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "参数不合法",
		})
		return
	}
	if req.ActionType != 1 && req.ActionType != 2 {
		c.JSON(http.StatusOK, serializer.CommentResponse{
			StatusCode: serializer.PermDenied,
			StatusMsg:  "参数不合法",
		})
		return
	}
	// 构造响应,返回评论内容
	userid, _ := c.Get("userid")
	resp := service.FavoritePostService(&req, userid.(int))
	c.JSON(http.StatusOK, resp)
}
func FavoriteList(c *gin.Context) {
	var req serializer.LikeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.CommListResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "参数不合法",
		})
		return
	}
	reqUserId, _ := c.Get("userid")
	i := reqUserId.(int)
	resp := service.FavoriteListService(&req, i)
	c.JSON(http.StatusOK, resp)
}
