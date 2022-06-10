package api

import (
	"basictiktok/serializer"
	"basictiktok/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户注册
func Register(c *gin.Context){
	var req serializer.RegisterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	resp := service.RegisterService(&req)
	c.JSON(http.StatusOK,resp)
}

// 用户登录
func Login(c *gin.Context){
	var req serializer.LoginRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	resp := service.LoginService(&req)
	c.JSON(http.StatusOK,resp)
}

// 查询用户信息
func QueryUserInfo(c *gin.Context)  {
	var req serializer.UserInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	resp := service.QueryUserInfoService(&req)
	c.JSON(http.StatusOK,resp)
}