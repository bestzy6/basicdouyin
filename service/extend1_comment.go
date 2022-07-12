package service

import (
	"basictiktok/dao"
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
	"time"
)

// CommentPostService 在对应的视频下添加评论
func CommentPostService(req *serializer.CommentRequest, userId int) *serializer.CommentResponse {
	var (
		resp    serializer.CommentResponse
		err     error
		postDao = dao.NewPostDaoInstance()
		userDao = dao.NewUserDaoInstance()
	)

	if req.ActionType == 1 {
		CommentIDGenerator, _ := util.NewGenerator(util.COMMENT)
		//创建评论
		post := model.Post{
			Id:         CommentIDGenerator.NextId(),
			VideoId:    int64(req.VideoId),
			UserId:     int64(userId),
			Content:    req.CommentText,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"), // 时间保持一致
		}
		err = postDao.Creat(&post)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = err.Error()
			resp.Comment = serializer.Comment{}
		}
		//
		user, _ := userDao.QueryUserByID(int64(userId))
		userTmp := serializer.User{
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			ID:            int64(userId),
			Name:          user.UserName,
		}
		//
		comment := serializer.Comment{
			Content:    post.Content,
			CreateDate: post.CreateTime,
			ID:         post.Id,
			User:       userTmp,
		}
		resp.Comment = comment
	} else {
		//删除评论
		err = postDao.Delete(req.CommentId, req.VideoId)
		if err != nil {
			resp.StatusCode = serializer.UnknownError
			resp.StatusMsg = err.Error()
		}
		resp.Comment = serializer.Comment{}
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "评论成功！"
	return &resp
}

func CommentListService(req *serializer.CommentListRequest) *serializer.CommListResponse {
	var (
		resp    serializer.CommListResponse
		postDao = dao.NewPostDaoInstance()
		userDao = dao.NewUserDaoInstance()
	)
	videoId := req.VideoId
	//查评论表
	commentList, err := postDao.QueryPostByVideoId(int64(videoId))
	if err != nil {
		util.Log().Error("查询失败:", err)
	}

	ansList := make([]*serializer.Comment, 0, len(commentList))
	for i := 0; i < len(commentList); i++ {
		//查用户表
		user, err := userDao.QueryUserByID(commentList[i].UserId)
		if err != nil {
			resp.StatusCode = serializer.OK
			resp.StatusMsg = err.Error()
			resp.CommentList = nil
			return &resp
		}
		//创建评论
		temp := &serializer.Comment{
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreateTime,
			ID:         commentList[i].Id,
			User:       toSeriUser(user),
		}
		ansList = append(ansList, temp)
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "评论列表查询成功"
	resp.CommentList = ansList
	return &resp
}

func toSeriUser(u model.User) serializer.User {
	return serializer.User{
		ID:            int64(u.ID),
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		Name:          u.UserName,
	}
}
