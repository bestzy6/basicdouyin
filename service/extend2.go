package service

import (
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/serializer"
)

// FollowService 关注服务（还需要添加对数据库的操作）
func FollowService(req *serializer.FollowRequest) *serializer.FollowResponse {
	var resp serializer.FollowResponse
	user := graphdb.User{ID: req.ReqUserId} //需要根据token修改
	targetUser := graphdb.User{ID: req.ToUserId}
	err := user.Follow(&targetUser)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	//以下操作异步修改数据库
	//
	return &resp
}

// UnFollowService 取消关注服务
func UnFollowService(req *serializer.FollowRequest) *serializer.FollowResponse {
	var resp serializer.FollowResponse
	user := graphdb.User{ID: req.ReqUserId} //需要根据token修改
	targetUser := graphdb.User{ID: req.ToUserId}
	err := user.UnFollow(&targetUser)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	//以下操作异步修改数据库
	//
	return &resp
}

// FollowersService 获取关注列表服务
func FollowersService(req *serializer.FollowListRequest) *serializer.FollowListResponse {
	var (
		resp  serializer.FollowListResponse
		users map[int]*graphdb.User
		err   error
	)
	user := graphdb.User{ID: req.UserId}
	if req.UserId == req.ReqUserId {
		users, err = user.MyFollowers()
	} else {
		reqUser := graphdb.User{ID: req.ReqUserId} //需要根据token修改
		users, err = user.Followers(&reqUser)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		resp.UserList = nil
		return &resp
	}
	//
	resp.UserList = toSeriaUsers(users)
	resp.StatusCode = serializer.OK
	resp.StatusMsg = ""
	return &resp
}

// FolloweesService 获取粉丝列表服务
func FolloweesService(req *serializer.FolloweesRequest) *serializer.FolloweesResponse {
	var (
		resp  serializer.FolloweesResponse
		users map[int]*graphdb.User
		err   error
	)
	user := graphdb.User{ID: req.UserId}
	if req.UserId == req.ReqUserId {
		users, err = user.MyFollowees()
	} else {
		reqUser := graphdb.User{ID: req.UserId} //需要根据token修改
		users, err = user.Followees(&reqUser)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		resp.UserList = nil
		return &resp
	}
	//
	resp.UserList = toSeriaUsers(users)
	resp.StatusCode = serializer.OK
	resp.StatusMsg = ""
	return &resp
}

// 将DTO的users转换为VO的serializer.Users
func toSeriaUsers(users map[int]*graphdb.User) []serializer.User {
	ansUsers := make([]serializer.User, 0, len(users))
	for _, v := range users {
		ansUser := serializer.User{
			ID:            int64(v.ID),
			Name:          v.Name,
			FollowCount:   int64(v.FollowCount),
			FollowerCount: int64(v.FollowerCount),
			IsFollow:      v.IsFollow,
		}
		ansUsers = append(ansUsers, ansUser)
	}
	return ansUsers
}

// 将graph.User转换为model.User
func toModelUser(user *graphdb.User) *model.User {
	return nil
}
