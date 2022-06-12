package api

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/service"
	"basictiktok/util"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//

func PublishVideo(c *gin.Context) {

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}

	// 获取视频以及图片存储目录
	video_dictory := os.Getenv("VIDEO_DICTORY")
	picture_dictory := os.Getenv("PICTURE_DICTORY")

	filename := filepath.Base(data.Filename)
	user := serializer.User{ID: 1}
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join(video_dictory, finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}

	// 将视频截取第一帧作为视频封面
	reader := util.ReadFrameAsJpeg(saveFile, 1)
	img, err := imaging.Decode(reader)
	if err != nil {
		fmt.Println(err.Error())
	}

	outputName := "out" + ".jpeg"
	outputName = fmt.Sprintf("%d_%s", user.ID, outputName)
	outputName = filepath.Join(picture_dictory, outputName)
	err = imaging.Save(img, outputName)
	if err != nil {
		fmt.Println(err.Error())

	}

	var video model.Video
	//userClaim, err := middleware.ParseToken(c.PostForm("token"))
	//if err != nil {
	//	c.JSON(http.StatusOK, serializer.ActionResponse{
	//		StatusCode: serializer.ParamInvalid,
	//		StatusMsg:  "请求参数错误",
	//	})
	//}
	// 生成视频信息
	//video.UserID = int64(userClaim.UserID)
	// 投稿用户id
	video.UserID = 1
	// 视频封面url
	video.CoverURL = outputName
	// 视频评论数
	video.CommentCount = 0
	// 视频点赞人数
	video.FavoriteCount = 0
	// 视频播放地址
	video.PlayURL = saveFile
	// 视频标题
	video.Title = c.PostForm("title")
	// 视频添加时间
	video.AddTime = time.Now()
	// 将视频信息插入数据库
	err = model.CreateAVideo(&video)
	if err != nil {
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

func VideoList(c *gin.Context) {
	var req serializer.FeedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, serializer.RegisterResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	resp := service.FindVideoBeforeTimeService(&req)
	c.JSON(http.StatusOK, resp)
}
