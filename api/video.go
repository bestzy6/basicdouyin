package api

import (
	"basictiktok/serializer"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

//

func PublishVideo(c *gin.Context) {
	//var req serializer.ActionRequest
	//if err := c.ShouldBindQuery(&req); err != nil {
	//	c.JSON(http.StatusOK, serializer.ActionResponse{
	//		StatusCode: serializer.ParamInvalid,
	//		StatusMsg:  "请求参数错误",
	//	})
	//	return
	//}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	fmt.Println(data)

	filename := filepath.Base(data.Filename)
	user := serializer.User{ID: 1}
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("../videodata/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}

	c.JSON(http.StatusOK, serializer.ActionResponse{
		StatusCode: serializer.OK,
		StatusMsg:  finalName + "upload successfully",
	})
}
