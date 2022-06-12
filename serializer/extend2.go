package serializer

// FollowResponse 关注操作
type FollowResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type FollowRequest struct {
	Token      string `form:"token" json:"token" binding:"required"`             //用户鉴权token
	ToUserId   int    `form:"to_user_id" json:"to_user_id" binding:"required"`   //对方用户id
	ActionType int    `form:"action_type" json:"action_type" binding:"required"` //1-关注，2-取消关注
	ReqUserId  int
}

// FollowListResponse 关注列表
type FollowListResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserList   []User `json:"user_list"`   // 用户信息列表
}

type FollowListRequest struct {
	Token     string `form:"token" json:"token" binding:"required"`     //用户鉴权token
	UserId    int    `form:"user_id" json:"user_id" binding:"required"` //用户id
	ReqUserId int
}

// FolloweesResponse 粉丝列表
type FolloweesResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserList   []User `json:"user_list"`   // 用户列表
}

type FolloweesRequest struct {
	Token     string `form:"token" json:"token" binding:"required"`     //用户鉴权token
	UserId    int    `form:"user_id" json:"user_id" binding:"required"` //用户id
	ReqUserId int
}
