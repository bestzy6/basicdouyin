package service

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
	"time"
)

// CommentPostService 在对应的视频下添加评论
func CommentPostService(req *serializer.CommentRequest) *serializer.CommentResponse {
	var resp serializer.CommentResponse
	var user serializer.User
	var comment serializer.Comment
	// 根据 token 获取userid ，根据userid 查询user 信息
	// userid := req.Token
	// user = getUser(userid)   和佳佳的表对接下
	if req.ActionType == 1 { // 添加评论（video_id -> comment表
		comment = serializer.Comment{
			Content:    req.CommentText,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			ID:         int64(req.CommentId),
			User:       user, // get 一波user
		}
		post := model.Post{
			Id:         int64(req.CommentId),
			VideoId:    int64(req.VideoId),
			UserId:     user.ID,
			Content:    req.CommentText,
			DiggCount:  10,                 // 每一条视频的评论总数
			CreateTime: comment.CreateDate, // 时间保持一致
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

func CommentListService(req *serializer.CommentListRequest) *serializer.CommListResponse {
	var resp serializer.CommListResponse
	videoId := req.VideoId
	//查表操作
	commentList, err := model.NewPostDaoInstance().QueryPostByVideoId(int64(videoId))
	if err != nil {
		util.Log().Error("查询失败:", err)
	}
	var commentTmp1 []*serializer.Comment
	for _, v := range commentList {
		commentTmp := serializer.Comment{
			Content:    v.Content,
			CreateDate: v.CreateTime,
			ID:         v.Id,
			User:       serializer.User{},
		}
		//commentTmp.User = v.UserId
		//commentTmp.CreateDate = v.CreateTime
		commentTmp1 = append(commentTmp1, &commentTmp)
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "评论列表查询成功"
	resp.CommentList = commentTmp1
	return &resp
}
