package service

import (
	"basictiktok/graphdb"
	"basictiktok/serializer"
)

// FollowService 关注服务（还需要添加对数据库的操作）
func FollowService(req *serializer.FollowRequest) *serializer.FollowResponse {
	var resp serializer.FollowResponse
	user := graphdb.User{ID: 1} //需要根据token修改
	targetUser := graphdb.User{ID: req.ToUserId}
	var err error
	if req.ActionType == 1 {
		err = user.Follow(&targetUser)
	} else {
		err = user.UnFollow(&targetUser)
	}
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	return &resp
}

// FollowersService 获取关注列表服务
func FollowersService(req *serializer.FollowListRequest) *serializer.FollowListResponse {
	var resp serializer.FollowListResponse
	reqUser := graphdb.User{ID: 1} //需要根据token修改
	user := graphdb.User{ID: req.UserId}
	users, err := user.Followers(&reqUser)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		resp.UserList = nil
		return &resp
	}
	//
	ansUsers := make([]serializer.User, len(users))
	for i := 0; i < len(users); i++ {
		ansUser := serializer.User{
			ID:            int64(users[i].ID),
			Name:          users[i].Name,
			FollowCount:   int64(users[i].FollowCount),
			FollowerCount: int64(users[i].FollowerCount),
			IsFollow:      users[i].IsFollow,
		}
		ansUsers[i] = ansUser
	}
	resp.UserList = ansUsers
	resp.StatusCode = serializer.OK
	resp.StatusMsg = ""
	return &resp
}

// FolloweesService 获取粉丝列表服务
func FolloweesService(req *serializer.FolloweesRequest) *serializer.FolloweesResponse {
	var resp serializer.FolloweesResponse
	reqUser := graphdb.User{ID: 1} //需要根据token修改
	user := graphdb.User{ID: req.UserId}
	users, err := user.Followees(&reqUser)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		resp.UserList = nil
		return &resp
	}
	//
	ansUsers := make([]serializer.User, len(users))
	for i := 0; i < len(users); i++ {
		ansUser := serializer.User{
			ID:            int64(users[i].ID),
			Name:          users[i].Name,
			FollowCount:   int64(users[i].FollowCount),
			FollowerCount: int64(users[i].FollowerCount),
			IsFollow:      users[i].IsFollow,
		}
		ansUsers[i] = ansUser
	}
	resp.UserList = ansUsers
	resp.StatusCode = serializer.OK
	resp.StatusMsg = ""
	return &resp
}
