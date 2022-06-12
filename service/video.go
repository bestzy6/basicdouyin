package service

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/server/middleware"
)

func PublishService(req *serializer.ActionRequest) {

}

func FindVideoBeforeTimeService(req *serializer.FeedRequest) *serializer.FeedResponse {

	var resp serializer.FeedResponse
	var videoList []serializer.Video

	if req.Token != "" {
		userClaim, err := middleware.ParseToken(req.Token)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = "未知错误"
			return &resp
		}

	}
	videos, err := model.FindVideoBeforeTime(req.LatestTime)
	for k := range videos {
		var user serializer.User
		userInfo, err := model.QueryUserByID(videos[k].UserID)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = "未知错误"
			return &resp
		}
		user.ID = int64(userInfo.ID)
		user.FollowerCount = userInfo.FollowerCount
		user.FollowCount = userInfo.FollowCount

	}
}
