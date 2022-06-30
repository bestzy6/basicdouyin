package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
)

// FavoritePostService 点赞操作
//func FavoritePostService(req *serializer.LikesRequest, userid int) *serializer.LikesResponse {
//	var resp serializer.LikesResponse
//	vid := req.VideoId
//	newV := model.NewVideoClDaoInstance()
//	newV.AddFavorite(int64(vid)) // 更新冗余表的点赞总数
//	num, _ := newV.QueryByVideoId(int64(vid))
//	// 更新video 表的评论总数字段
//	// 调用下佳佳更新video表就完事了         这样一个一个更新影响效率
//	videoDao := model.NewVideoDaoInstance()
//	videoDao.AddFavorite(int64(vid))
//
//	if req.ActionType == 1 { //点赞操作
//		fPost := model.FavoritePost{
//			UserId:    int64(userid), // 根据token 获得 user_id
//			VideoId:   int64(req.VideoId),
//			DiggCount: int32(num.FavoriteCount),
//		}
//		if err := model.NewFavoritePostDaoInstance().CreateFPost(&fPost); err != nil {
//			util.Log().Error("点赞失败:", err)
//		}
//	} else {
//		if err := newV.DeFavorite(int64(vid)); err != nil { // 根据给定的条件更新单个属性
//			util.Log().Error("取消点赞失败:", err)
//		} else {
//			videoDao.DeleteFavorite(int64(vid))
//		}
//	}
//	resp.StatusCode = serializer.OK
//	resp.StatusMsg = "点赞成功"
//	return &resp
//}

func FavoritePostService(req *serializer.LikesRequest, userid int) *serializer.LikesResponse {
	var (
		resp serializer.LikesResponse
		err  error
	)
	user := graphdb.User{ID: userid}
	vedio := graphdb.Video{ID: req.VideoId}
	if req.ActionType == 1 {
		//点赞
		err = user.Favorite(&vedio)
	} else {
		//取消点赞
		err = user.UnFavorite(&vedio)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = err.Error()
	}

	resp.StatusCode = serializer.OK
	resp.StatusMsg = "点赞成功"
	return &resp
}

// FavoriteListService 获取点赞列表
func FavoriteListService(req *serializer.LikeListRequest, myUserId int) *serializer.LikeListResponse {
	var resp serializer.LikeListResponse
	userId := req.UserId
	favoritePostDao := model.NewFavoritePostDaoInstance()
	videoPost, err := favoritePostDao.QueryFavoritePostById(int64(userId))
	if err != nil {
		util.Log().Error("获取点赞列表失败:", err)
	}
	videoLs := favoritePostDao.GetVideoIdList(videoPost)
	results, err1 := favoritePostDao.QueryPostByVedioId(videoLs) // 拿到所有相关的video 实体
	if err1 != nil {
		util.Log().Error("获取点赞列表失败:", err)
	}
	var videoTmpIndex []*serializer.Video

	for _, result := range results {
		userTmp, _ := model.QueryUserByID(result.UserID)
		user := serializer.User{
			FollowCount:   userTmp.FollowCount,
			FollowerCount: userTmp.FollowerCount,
			ID:            int64(userTmp.ID),
			Name:          userTmp.UserName,
		}
		user.IsFollow, _ = graphdb.IsFollow(myUserId, userTmp.ID)
		// 通过token获取你的userid，判断你是否关注这个视频作者
		videoTmp := serializer.Video{
			Author:        user,
			CommentCount:  result.CommentCount,
			CoverURL:      result.CoverURL,
			FavoriteCount: result.FavoriteCount,
			ID:            result.ID,
			IsFavorite:    false,
			PlayURL:       result.PlayURL,
			Title:         result.Title,
		}
		videoTmpIndex = append(videoTmpIndex, &videoTmp)
	}

	resp.StatusCode = serializer.OK
	resp.StatusMsg = "查询点赞列表成功"
	resp.VideoList = videoTmpIndex
	return &resp

}

// IsFavorite 获取用户点赞的视频
func IsFavorite(userId, videoId int64) bool {
	newFPO := model.NewFavoritePostDaoInstance()
	tmp, err := newFPO.QueryFavoritePostById(userId)
	if err != nil {
		util.Log().Error("获取点赞视频失败:", err)
	}
	videoLs := newFPO.GetVideoIdList(tmp)
	for _, v := range videoLs {
		if v == videoId {
			return true
		}
	}
	return false
}
