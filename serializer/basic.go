package serializer

import (
	"mime/multipart"
)

// ListResponse 发布列表
type ListResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户发布的视频列表
}

type ListRequest struct {
	Token  string `form:"token" json:"token" binding:"required"` //用户鉴权token
	UserId int64  `form:"user_id" json:"user_id" binding:"required"`
}

// ActionResponse 投稿接口
type ActionResponse struct {
	StatusCode ErrNo  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type ActionRequest struct {
	Token string                `form:"token" json:"token" binding:"required"` //用户鉴权token
	Title string                `form:"title" json:"title" binding:"required"`
	Data  *multipart.FileHeader `form:"data" json:"data" binding:"required"`
}

// UserInfoResponse 用户信息返回
type UserInfoResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	User       User   `json:"user"`        // 用户信息
}

// UserInfoRequest 用户信息输入
type UserInfoRequest struct {
	Token  string `form:"token" json:"token" binding:"required"`     //用户鉴权token
	UserId int64  `form:"user_id" json:"user_id" binding:"required"` //用户id
}

// LoginResponse 用户登入响应
type LoginResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// LoginRequest 用户登入输入
type LoginRequest struct {
	UserName       string `form:"user_name" json:"user_name" binding:"required"`             // 用户名
	PasswordDigest string `form:"password_digest" json:"password_digest" binding:"required"` // 用户密码
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// RegisterRequest 用户注册输入
type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required"` // 用户名
	Password string `form:"password" json:"password" binding:"required"` // 用户密码
	Nickname string `form:"nickname" json:"nickname"`                    // 用户昵称
}

// FeedResponse 视频流接口
type FeedResponse struct {
	NextTime   int64   `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 视频列表
}

// FeedRequset 视频流请求
type FeedRequest struct {
	LatestTime int64  `form:"lastet_time" json:"latest_time"` // 返回当前指定时间之前上传的视频视频
	Token      string `form:"token" json:"token"`             //用户鉴权token
}
