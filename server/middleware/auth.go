package middleware

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exist := c.GetQuery("token")
		if !exist {
			//返回json可能需要修改（2022年6月8日10点47分，zy）
			c.JSON(http.StatusOK, gin.H{
				"msg": "无token",
			})
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		//解析出userid放置到上下文中，将数据向后传递。
		c.Set("userid", claims.UserID)
	}
}

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid != nil {
			user, err := model.GetUser(uid)
			if err == nil {
				c.Set("user", &user)
			}
		}
		c.Next()
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user"); user != nil {
			if _, ok := user.(*model.User); ok {
				c.Next()
				return
			}
		}

		c.JSON(200, serializer.CheckLogin())
		c.Abort()
	}
}
