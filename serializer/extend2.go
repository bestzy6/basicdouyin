package serializer

// FollowResponse 关注操作
type FollowResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// FollowListResponse 关注列表
type FollowListResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	UserList   []User  `json:"user_list"`   // 用户信息列表
}

// FolloweesResponse 粉丝列表
type FolloweesResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	UserList   []User  `json:"user_list"`   // 用户列表
}
