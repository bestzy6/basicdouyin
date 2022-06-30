package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FindVideoBeforeTimeService(req *serializer.FeedRequest, userid int) *serializer.FeedResponse {
	var resp serializer.FeedResponse
	var user model.User
	user.ID = userid
	// 找出请求时间之前的30条视频信息返回
	videos, err := model.FindVideoBeforeTime(req.LatestTime)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	if len(videos) == 0 {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "服务器没有视频"
		return &resp
	}
	videoList := make([]serializer.Video, 0, 30)
	// 将视频信息与上传作者信息进行绑定返回
	for k := range videos {
		upInfo, err := model.QueryUserByID(videos[k].UserID)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = "未知错误"
			return &resp
		}
		uper := serializer.User{
			ID:            int64(upInfo.ID),
			FollowerCount: upInfo.FollowerCount,
			FollowCount:   upInfo.FollowCount,
		}

		var videoRes serializer.Video
		if user.ID != 0 {
			// 判断当前登录用户是否关注当前视频作者
			isFollow, err := graphdb.IsFollow(user.ID, int(uper.ID))
			if err != nil {
				resp.StatusCode = serializer.UnknownError
				resp.StatusMsg = "未知错误"
				return &resp
			}
			uper.IsFollow = isFollow
			// 判断是否点赞
			videoRes.IsFavorite = graphdb.IsFavorite(user.ID, int(videos[k].ID))
		} else {
			videoRes.IsFavorite = false
		}
		uper.Name = upInfo.UserName
		videoRes.Author = uper
		videoRes.Title = videos[k].Title
		videoRes.PlayURL = videos[k].PlayURL
		videoRes.CoverURL = videos[k].CoverURL
		videoRes.FavoriteCount = videos[k].FavoriteCount
		videoRes.CommentCount = videos[k].CommentCount
		videoRes.ID = videos[k].ID
		//
		videoList = append(videoList, videoRes)
	}
	resp.VideoList = videoList
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	// 返回结果已经进行排序，因此第0个视频发布时间即下次最新时间
	resp.NextTime = videos[0].AddTime
	return &resp
}

func ListVideosService(req *serializer.ListRequest, userid int) *serializer.ListResponse {
	var resp serializer.ListResponse
	user, err := model.QueryUserByID(req.UserId)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	userRes := serializer.User{
		ID:            int64(user.ID),
		Name:          user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
	}
	userRes.IsFollow, err = graphdb.IsFollow(userid, user.ID)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}

	videos, err := model.QueryVideoListByUserID(user.ID)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	videoList := make([]serializer.Video, 0, len(videos))
	for k := range videos {
		video := serializer.Video{
			Author:        userRes,
			IsFavorite:    graphdb.IsFavorite(user.ID, int(videos[k].ID)),
			Title:         videos[k].Title,
			PlayURL:       videos[k].PlayURL,
			CoverURL:      videos[k].CoverURL,
			FavoriteCount: videos[k].FavoriteCount,
			CommentCount:  videos[k].CommentCount,
			ID:            videos[k].ID,
		}
		videoList = append(videoList, video)
	}
	resp.VideoList = videoList
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	return &resp
}

func ActionService(req *serializer.ActionRequest, userid int, host string) *serializer.ActionResponse {
	prefix := getFileName(userid)
	split := strings.Split(filepath.Base(req.Data.Filename), ".")
	vedioFileName := prefix + "." + split[1] //视频文件名
	saveVedioPath := filepath.Join(util.VEDIO, vedioFileName)
	err := saveUploadedFile(req.Data, saveVedioPath)
	if err != nil {
		util.Log().Error("保存文件出错！\n", err)
		return &serializer.ActionResponse{
			StatusCode: serializer.UnknownError,
			StatusMsg:  "保存文件出错！",
		}
	}
	// 将视频截取第一帧作为视频封面
	reader := util.ReadFrameAsJpeg(saveVedioPath, 1)
	img, err := imaging.Decode(reader)
	if err != nil {
		util.Log().Error("截取视频封面错误！\n", err)
		return &serializer.ActionResponse{
			StatusCode: serializer.UnknownError,
			StatusMsg:  "截取视频封面错误！",
		}
	}
	imgFileName := prefix + ".jpeg" //截图文件名
	saveImgPath := filepath.Join(util.IMG, imgFileName)
	err = imaging.Save(img, saveImgPath)
	if err != nil {
		util.Log().Error("保存封面错误！\n", err)
		return &serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "保存封面错误！",
		}
	}
	video := model.Video{
		UserID:        int64(userid),
		CoverURL:      "http://" + host + "/static/img/" + imgFileName,
		CommentCount:  0,
		FavoriteCount: 0,
		PlayURL:       "http://" + host + "/static/video/" + vedioFileName,
		Title:         req.Title,
		AddTime:       time.Now().Unix(),
	}
	err = video.Create()
	if err != nil {
		return &serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		}
	}
	graphVideo := toGraphVideo(&video)
	err = graphVideo.Create()
	if err != nil {
		fmt.Println(err)
		return &serializer.ActionResponse{
			StatusCode: serializer.ParamInvalid,
			StatusMsg:  "请求参数错误",
		}
	}
	return &serializer.ActionResponse{
		StatusCode: serializer.OK,
		StatusMsg:  "视频上传成功",
	}
}

func getFileName(userid int) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(userid))
	builder.WriteString("_")
	builder.WriteString(util.RandStringRunes(3))
	builder.WriteString("_")
	builder.WriteString(time.Now().Format("20060102150405"))
	return builder.String()
}

func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func toGraphVideo(v *model.Video) *graphdb.Video {
	return &graphdb.Video{
		ID:            int(v.ID),
		UserID:        int(v.UserID),
		CoverURL:      v.CoverURL,
		FavoriteCount: int(v.FavoriteCount),
	}
}
