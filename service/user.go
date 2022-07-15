package service

import (
	"basictiktok/dao"
	"basictiktok/graphdb"
	"basictiktok/model"
	"basictiktok/serializer"
	"basictiktok/util"
)

// RegisterService 用户注册
func RegisterService(req *serializer.RegisterRequest) *serializer.RegisterResponse {
	var resp serializer.RegisterResponse
	userIdGenerator, _ := util.NewGenerator(util.USERID)
	user := model.User{
		ID:       int(userIdGenerator.NextId()),
		UserName: req.Username,
		Password: util.PasswordWithMD5(req.Password),
		Nickname: req.Nickname,
	}
	userDao := dao.NewUserDaoInstance()
	// 通过注册信息新建一个用户
	if err := userDao.Create(&user); err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}

	// 返回用户注册消息
	token, err := util.CreateToken(user.ID)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	//添加信息至图数据库
	userGraphDao := graphdb.NewUserGraphDao()
	graphUser := model2graph(&user)
	err = userGraphDao.Create(graphUser)
	if err != nil {
		util.Log().Error("创建图用户失败！", err)
	}
	//
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	resp.Token = token
	resp.UserID = int64(user.ID)
	return &resp
}

// LoginService 用户登录
func LoginService(req *serializer.LoginRequest) *serializer.LoginResponse {
	var resp serializer.LoginResponse
	userDao := dao.NewUserDaoInstance()
	user, err := userDao.QueryUserByName(req.UserName)
	if err != nil {
		resp.StatusCode = serializer.UserNotExisted
		resp.StatusMsg = "用户名错误"
		return &resp
	}

	passwordWithMD5 := util.PasswordWithMD5(req.Password)
	if passwordWithMD5 != user.Password {
		resp.StatusCode = serializer.WrongPassword
		resp.StatusMsg = "密码错误"
		return &resp
	}

	token, err := util.CreateToken(user.ID)
	if err != nil {
		resp.StatusCode = serializer.UnknownError
		resp.StatusMsg = "未知错误"
		return &resp
	}
	resp.StatusCode = serializer.OK
	resp.StatusMsg = "ok"
	resp.Token = token
	resp.UserID = int64(user.ID)
	return &resp

}

// QueryUserInfoService 用户查询
func QueryUserInfoService(req *serializer.UserInfoRequest, userid int) *serializer.UserInfoResponse {
	var resp serializer.UserInfoResponse
	userDao := dao.NewUserDaoInstance()
	user, err := userDao.QueryUserByID(req.UserId)
	userGraphDao := graphdb.NewUserGraphDao()
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
	resp.User.IsFollow, _ = userGraphDao.HasFollow(userid, int(req.UserId))
	return &resp
}

// model的用户转换为graphdb的用户
func model2graph(user *model.User) *graphdb.User {
	toUser := &graphdb.User{
		ID:   user.ID,
		Name: user.UserName,
	}
	return toUser
}
