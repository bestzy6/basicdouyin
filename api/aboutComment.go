package api

import (
	"basictiktok/serializer"
	"basictiktok/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommentPost 状态检查页面，实现具体的逻辑，构造comment
func CommentPost(c *gin.Context) {
	var req serializer.CommentRequest
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
	//reqUserId, _ := c.Get("userid")
	resp := service.CommentPostService(&req)
	c.JSON(http.StatusOK, resp)
}

// CommentList 获取评论列表
func CommentList(c *gin.Context) {
	var req serializer.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.CommListResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "参数不合法",
		})
		return
	}
	resp := service.CommentListService(&req)
	c.JSON(http.StatusOK, resp)
}
