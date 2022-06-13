package middleware

import (
	"basictiktok/serializer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthToken token验证
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exist := c.GetQuery("token")
		if !exist {
			token, exist = c.GetPostForm("token")
			if !exist {
				c.JSON(http.StatusOK, gin.H{
					"status_code": serializer.ParamInvalid,
					"status_msg":  "参数错误！",
				})
				c.Abort()
				return
			}
		}
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": serializer.PermDenied,
				"status_msg":  "无操作权限",
			})
			c.Abort()
			return
		}
		c.Set("userid", claims.UserID)
		c.Next()
	}
}

//// CurrentUser 获取登录用户
//func CurrentUser() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		session := sessions.Default(c)
//		uid := session.Get("user_id")
//		if uid != nil {
//			user, err := model.GetUser(uid)
//			if err == nil {
//				c.Set("user", &user)
//			}
//		}
//		c.Next()
//	}
//}
//
//// AuthRequired 需要登录
//func AuthRequired() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if user, _ := c.Get("user"); user != nil {
//			if _, ok := user.(*model.User); ok {
//				c.Next()
//				return
//			}
//		}
//
//		//c.JSON(200, serializer.CheckLogin())
//		c.Abort()
//	}
//}
