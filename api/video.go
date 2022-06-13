package api

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/server/middleware"
	"basictiktok/service"
	"basictiktok/util"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/u2takey/go-utils/uuid"
	"net/http"
	"path/filepath"
	"strings"
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
	video_dictory := "./static/video/"
	picture_dictory := "./static/img/"

	filename := filepath.Base(data.Filename)
	// 通过token获取登录用户id
	id, _ := c.Get("userid")
	userid := id.(int)
	finalName := fmt.Sprintf("%d_%s", userid, filename)
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
		fmt.Println("decode", err.Error())
		c.JSON(http.StatusOK, serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	// 通过UUID生成唯一的视频封面名
	name := uuid.NewUUID()
	name = strings.Replace(name, "-", "", -1)
	outputName := name + ".jpeg"
	outputName = fmt.Sprintf("%d_%s", userid, outputName)
	savePicture := filepath.Join(picture_dictory, outputName)
	err = imaging.Save(img, savePicture)
	if err != nil {
		fmt.Println(err.Error())

	}
	video := model.Video{
		UserID:        int64(userid),
		CoverURL:      "http://" + c.Request.Host + "/static/img/" + outputName,
		CommentCount:  0,
		FavoriteCount: 0,
		PlayURL:       "http://" + c.Request.Host + "/static/video/" + finalName,
		Title:         c.PostForm("title"),
	}
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

// 视频流接口
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

// 用户视频发布列表
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
