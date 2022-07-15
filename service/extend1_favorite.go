package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/mq"
	"basictiktok/serializer"
)

func FavoritePostService(req *serializer.LikesRequest, userid int) *serializer.LikesResponse {
	var (
		resp serializer.LikesResponse
		err  error
	)
	video := graphdb.Video{ID: req.VideoId}
	userGraphDao := graphdb.NewUserGraphDao()
	if req.ActionType == 1 {
		//点赞
		err = userGraphDao.Favorite(userid, req.VideoId)
	} else {
		//取消点赞
		err = userGraphDao.UnFavorite(userid, req.VideoId)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = err.Error()
	}
	//mysql异步更新
	msg := &mq.UserMessage{
		ToVideo: videoG2M(&video),
	}
	if req.ActionType == 1 {
		msg.OpNum = mq.Favorite
	} else {
		msg.OpNum = mq.UnFavorite
	}
	mq.ToModelUserMQ <- msg
	//
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "点赞成功"
	return &resp
}

// FavoriteListService 获取点赞列表
func FavoriteListService(req *serializer.LikeListRequest, myUserId int) *serializer.LikeListResponse {
	var (
		resp serializer.LikeListResponse
		list map[int]*graphdb.Video
		err  error
	)
	userGraphDao := graphdb.NewUserGraphDao()
	if req.UserId == myUserId {
		list, err = userGraphDao.MyFavoriteList(myUserId)
	} else {
		list, err = userGraphDao.FavoriteList(req.UserId, myUserId)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = err.Error()
		resp.VideoList = nil
	}
	videoList := make([]*serializer.Video, 0, len(list))
	for _, v := range list {
		ansVideo := &serializer.Video{
			ID:            int64(v.ID),
			CoverURL:      v.CoverURL,
			FavoriteCount: int64(v.FavoriteCount),
		}
		videoList = append(videoList, ansVideo)
	}
	//
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "获取点赞列表成功"
	resp.VideoList = videoList
	return &resp
}

func videoG2M(m *graphdb.Video) *model.Video {
	return &model.Video{ID: int64(m.ID)}
}
