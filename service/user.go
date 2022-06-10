package service

import (
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
)

// RegisterService 用户注册
func RegisterService(req *serializer.RegisterRequest) *serializer.RegisterResponse {
	var resp serializer.RegisterResponse
	var user model.User
	user.UserName = req.UserName
	// 获取md5加密后的密码
	user.PasswordDigest = util.PasswordWithMD5(req.PasswordDigest)
	user.Nickname = req.Nickname

	// 通过注册信息新建一个用户
	if err := model.CreateAUser(&user); err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}

	// 通过用户名查询新建用户的信息
	newUser, err := model.QueryAUser(user.UserName)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}

	// 返回用户注册消息
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	resp.Token = "token"
	resp.UserID = int64(newUser.ID)
	return &resp
}

// LoginService 用户登录
func LoginService(req *serializer.LoginRequest) *serializer.LoginResponse {

	var resp serializer.LoginResponse
	user, err := model.QueryAUser(req.UserName)
	if err != nil {
		resp.StatusCode = serializer.UserNotExisted
		resp.StatusMsg = "用户名错误"
		return &resp
	}

	passwordWithMD5 := util.PasswordWithMD5(req.PasswordDigest)
	if passwordWithMD5 != user.PasswordDigest {
		resp.StatusCode = serializer.WrongPassword
		resp.StatusMsg = "密码错误"
		return &resp
	}

	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	resp.Token = "token"
	resp.UserID = int64(user.ID)
	return &resp

}

// QueryUserInfoService 用户查询
func QueryUserInfoService(req *serializer.UserInfoRequest) *serializer.UserInfoResponse {
	var resp serializer.UserInfoResponse
	user, err := model.QueryUser(int64(req.UserId))
	if err != nil {
		resp.StatusCode = serializer.ParamInvalid
		resp.StatusMsg = "用户id错误"
		return &resp
	}
	// 返回用户查询信息
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	resp.User.ID = int64(user.ID)
	resp.User.Name = user.UserName
	resp.User.FollowCount = user.FollowCount
	resp.User.FollowerCount = user.FollowerCount
	// 这个需要查询用户关注信息
	resp.User.IsFollow = true
	return &resp
}
