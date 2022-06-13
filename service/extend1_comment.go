package service

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
	"time"
)

// CommentPostService 在对应的视频下添加评论
func CommentPostService(req *serializer.CommentRequest, userId int) *serializer.CommentResponse {
	var resp serializer.CommentResponse
	var comment serializer.Comment
	user, _ := model.QueryUserByID(int64(userId))
	userTmp := serializer.User{
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		ID:            int64(userId),
		IsFollow:      false,
		Name:          user.UserName,
	}
	// 获取评论数
	newV := model.NewVideoClDaoInstance()
	newV.AddComment(int64(req.VideoId))               // 先更新冗余表
	num, _ := newV.QueryByVideoId(int64(req.VideoId)) // 拿到最新的评论数
	if req.ActionType == 1 {                          // 添加评论
		comment = serializer.Comment{
			Content:    req.CommentText,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			ID:         int64(req.CommentId),
			User:       userTmp, // get 一波user
		}
		post := model.Post{
			Id:         int64(req.CommentId),
			VideoId:    int64(req.VideoId),
			UserId:     int64(userId),
			Content:    req.CommentText,
			DiggCount:  int32(num.CommentCount), // 最新值，之后再改成消息队列的形式
			CreateTime: comment.CreateDate,      // 时间保持一致
		}
		if err := model.NewPostDaoInstance().CreatePost(&post); err != nil {
			util.Log().Error("添加评论失败:", err)
		}
	} else {
		// 删除评论（查表）
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "评论成功"
	resp.Comment = comment
	return &resp
}

func CommentListService(req *serializer.CommentListRequest, userId int) *serializer.CommListResponse {
	var resp serializer.CommListResponse
	videoId := req.VideoId
	//查表操作
	commentList, err := model.NewPostDaoInstance().QueryPostByVideoId(int64(videoId))
	if err != nil {
		util.Log().Error("查询失败:", err)
	}
	user, _ := model.GetUser(userId)
	tmpUser := serializer.User{
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		ID:            int64(user.ID),
		IsFollow:      false,
		Name:          user.UserName,
	}
	var commentTmp1 []*serializer.Comment
	for _, v := range commentList {
		commentTmp := serializer.Comment{
			Content:    v.Content,
			CreateDate: v.CreateTime,
			ID:         v.Id,
			User:       tmpUser,
		}
		commentTmp1 = append(commentTmp1, &commentTmp)
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "评论列表查询成功"
	resp.CommentList = commentTmp1
	return &resp
}