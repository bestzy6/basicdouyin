package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/mq"
	"basictiktok/serializer"
	"strconv"
)

func FavoritePostService(req *serializer.LikesRequest, userid int) *serializer.LikesResponse {
	var (
		resp serializer.LikesResponse
		err  error
	)
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
	msg := strconv.Itoa(userid) + "_" + strconv.Itoa(req.VideoId) + "_" + strconv.Itoa(req.ActionType)
	mq.FavoriteProducerMsg <- msg
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
