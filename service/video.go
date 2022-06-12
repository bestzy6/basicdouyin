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
	var user model.User
	if req.Token != "" {
		userClaim, err := middleware.ParseToken(req.Token)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = "未知错误"
			return &resp
		}
		user.ID = userClaim.UserID

	}
	videos, err := model.FindVideoBeforeTime(req.LatestTime)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	for k := range videos {
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
			if IsFollow(user.ID, int(uper.ID)) {
				uper.IsFollow = true
			}
		}
	}
}

func IsFollow(user_id_1, user_id_2 int) bool {
	return true
}
