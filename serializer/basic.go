package serializer

// ListResponse 发布列表
type ListResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户发布的视频列表
}

// ActionResponse 投稿接口
type ActionResponse struct {
	StatusCode ErrNo   `json:"status_code"`
	StatusMsg  *string `json:"status_msg"`
}

// UserInfoResponse 用户信息
type UserInfoResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	User       User    `json:"user"`        // 用户信息
}

// LoginResponse 用户登入
type LoginResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	Token      *string `json:"token"`       // 用户鉴权token
	UserID     *int64  `json:"user_id"`     // 用户id
}

// RegisterResponse 用户注册
type RegisterResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// FeedResponse 视频流接口
type FeedResponse struct {
	NextTime   *int64  `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 视频列表
}
