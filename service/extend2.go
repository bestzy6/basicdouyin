package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
)

// FollowService 关注服务（还需要添加对数据库的操作）
func FollowService(req *serializer.FollowRequest) *serializer.FollowResponse {
	var resp serializer.FollowResponse
	user := graphdb.User{ID: req.ReqUserId} //需要根据token修改
	targetUser := graphdb.User{ID: req.ToUserId}
	//判断user和targetuser是否为同一人
	if user.ID == targetUser.ID {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "不能关注或取消关注自己！"
		return &resp
	}
	var err error
	if req.ActionType == 1 {
		err = user.Follow(&targetUser)
	} else {
		err = user.UnFollow(&targetUser)
	}
	if err != nil {
		util.Log().Error("neo4j关注错误\n", err)
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	modelUser := graph2model(&user)
	modelToUser := graph2model(&targetUser)
	if req.ActionType == 1 {
		err = modelUser.Follow(modelToUser)
	} else {
		err = modelUser.UnFollow(modelToUser)
	}
	if err != nil {
		util.Log().Error("mysql关注错误\n", err)
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	return &resp
}

// FollowersService 获取关注列表服务
func FollowersService(req *serializer.FollowListRequest) *serializer.FollowListResponse {
	var resp serializer.FollowListResponse
	reqUser := graphdb.User{ID: req.ReqUserId} //需要根据token修改
	user := graphdb.User{ID: req.UserId}
	//
	var (
		users map[int]*graphdb.User
		err   error
	)
	if reqUser.ID == user.ID {
		users, err = user.MyFollowers()
	} else {
		users, err = user.Followers(&reqUser)
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
	reqUser := graphdb.User{ID: req.ReqUserId} //需要根据token修改
	user := graphdb.User{ID: req.UserId}
	users, err := user.Followees(&reqUser)
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
