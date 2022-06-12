package service

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
)

// FavoritePostService 点赞操作
func FavoritePostService(req *serializer.LikesRequest) *serializer.LikesResponse {
	// 先验证token
	var resp serializer.LikesResponse
	vid := req.VideoId
	commentNums, _ := model.NewVideoClDaoInstance().QueryByVideoId(int64(vid))
	//vidObj, _ := model.QueryByVideoId(int64(vid))  等佳佳写完对接下
	//vidObj.FavoriteCount = commentNums.FavoriteCount
	// 更新表的点赞总数
	if req.ActionType == 1 { //点赞操作
		fPost := model.FavoritePost{
			UserId:    1, // 根据token 获得 user_id
			VideoId:   int64(req.VideoId),
			DiggCount: int32(commentNums.FavoriteCount + 1), // 这个怎么++呢 ，对应视频的点赞数++，根据video_id 找到对应的视频，然后把那个属性更新
		}
		if err := model.NewFavoritePostDaoInstance().CreateFPost(&fPost); err != nil {
			util.Log().Error("点赞失败:", err)
		}
	} else { //取消点赞

	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "点赞成功"
	return &resp
}

// FavoriteListService 获取点赞列表
func FavoriteListService(req *serializer.LikeListRequest) *serializer.LikeListResponse {
	// 验证token
	var resp serializer.LikeListResponse
	userId := req.UserId
	favoritePostDao := model.NewFavoritePostDaoInstance()
	videoPost, err := favoritePostDao.QueryFavoritePostById(int64(userId))
	if err != nil {
		util.Log().Error("点赞失败:", err)
	}
	videoLs := favoritePostDao.GetVideoIdList(videoPost)
	results, err1 := favoritePostDao.QueryPostByUserId(videoLs) // 拿到所有相关的video 实体
	if err1 != nil {
		util.Log().Error("点赞失败:", err)
	}
	var videoTmpIndex []*serializer.Video
	for _, result := range results {
		videoTmp := serializer.Video{
			Author:        serializer.User{},
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
