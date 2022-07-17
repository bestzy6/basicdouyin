package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/mq"
	"basictiktok/serializer"
	"basictiktok/util"
	"strconv"
)

// FollowService 关注服务（还需要添加对数据库的操作）
func FollowService(req *serializer.FollowRequest) *serializer.FollowResponse {
	var resp serializer.FollowResponse
	//判断user和targetuser是否为同一人
	if req.ReqUserId == req.ToUserId {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "不能关注或取消关注自己！"
		return &resp
	}
	userGraphDao := graphdb.NewUserGraphDao()
	var err error
	if req.ActionType == 1 {
		err = userGraphDao.Follow(req.ReqUserId, req.ToUserId)
	} else {
		err = userGraphDao.UnFollow(req.ReqUserId, req.ToUserId)
	}
	if err != nil {
		util.Log().Error("neo4j关注错误\n", err)
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	//mysql异步更新
	msg := strconv.Itoa(req.ReqUserId) + "_" + strconv.Itoa(req.ToUserId) + "_" + strconv.Itoa(req.ActionType)
	mq.FollowProducerMsg <- msg
	//
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	return &resp
}

// FollowersService 获取关注列表服务
func FollowersService(req *serializer.FollowListRequest) *serializer.FollowListResponse {
	var resp serializer.FollowListResponse
	var (
		users        map[int]*graphdb.User
		err          error
		userGraphDao = graphdb.NewUserGraphDao()
	)
	if req.ReqUserId == req.UserId {
		users, err = userGraphDao.MyFollowers(req.ReqUserId)
	} else {
		users, err = userGraphDao.Followers(req.UserId, req.ReqUserId)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		resp.UserList = nil
		return &resp
	}
	//
	ansUsers := make([]serializer.User, 0, len(users))
	for _, v := range users {
		ansUser := serializer.User{
			ID:       int64(v.ID),
			Name:     v.Name,
			IsFollow: v.IsFollow,
		}
		ansUsers = append(ansUsers, ansUser)
	}
	resp.UserList = ansUsers
	resp.StatusCode = serializer.OK
	resp.StatusMsg = ""
	return &resp
}

// FolloweesService 获取粉丝列表服务
func FolloweesService(req *serializer.FolloweesRequest) *serializer.FolloweesResponse {
	var resp serializer.FolloweesResponse
	userGraphDao := graphdb.NewUserGraphDao()
	users, err := userGraphDao.Followees(req.UserId, req.ReqUserId)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		resp.UserList = nil
		return &resp
	}
	//
	ansUsers := make([]serializer.User, 0, len(users))
	for _, v := range users {
		ansUser := serializer.User{
			ID:       int64(v.ID),
			Name:     v.Name,
			IsFollow: v.IsFollow,
		}
		ansUsers = append(ansUsers, ansUser)
	}
	resp.UserList = ansUsers
	resp.StatusCode = serializer.OK
	resp.StatusMsg = ""
	return &resp
}

func graph2model(user *graphdb.User) *model.User {
	toUser := &model.User{
		ID:       user.ID,
		UserName: user.Name,
	}
	return toUser
}
