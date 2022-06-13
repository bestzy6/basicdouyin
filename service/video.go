package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/serializer"
)

func FindVideoBeforeTimeService(req *serializer.FeedRequest, userid int) *serializer.FeedResponse {

	var resp serializer.FeedResponse
	var videoList []serializer.Video
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
	// 将视频信息与上传作者信息进行绑定返回
	for k := range videos {
		var videoRes serializer.Video
		var uper serializer.User
		upInfo, err := model.QueryUserByID(videos[k].UserID)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = "未知错误"
			return &resp
		}
		uper.ID = int64(upInfo.ID)
		uper.FollowerCount = upInfo.FollowerCount
		uper.FollowCount = upInfo.FollowCount
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
			videoRes.IsFavorite = IsFavorite(int64(user.ID), videos[k].ID)
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
	var userRes serializer.User
	var videoList []serializer.Video
	userRes.ID = int64(user.ID)
	userRes.Name = user.UserName
	userRes.IsFollow, err = graphdb.IsFollow(userid, user.ID)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	userRes.FollowCount = user.FollowCount
	userRes.FollowerCount = user.FollowerCount

	videos, err := model.QueryVideoListByUserID(user.ID)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	for k := range videos {
		var videoRes serializer.Video
		videoRes.Author = userRes
		// 判断是否点赞
		videoRes.IsFavorite = IsFavorite(int64(user.ID), videos[k].ID)
		videoRes.Title = videos[k].Title
		videoRes.PlayURL = videos[k].PlayURL
		videoRes.CoverURL = videos[k].CoverURL
		videoRes.FavoriteCount = videos[k].FavoriteCount
		videoRes.CommentCount = videos[k].CommentCount
		videoRes.ID = videos[k].ID
		videoList = append(videoList, videoRes)
	}

	resp.VideoList = videoList
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	return &resp
}
