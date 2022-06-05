package serializer

// LikesResponse 赞操作
type LikesResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// LikeListResponse 点赞列表
type LikeListResponse struct {
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户点赞视频列表
}

// CommentResponse 评论操作
type CommentResponse struct {
	Comment    Comment `json:"comment"`     // 评论成功返回评论内容，不需要重新拉取整个列表
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
}

// CommListResponse 评论列表
type CommListResponse struct {
	CommentList []Comment `json:"comment_list"` // 评论列表
	StatusCode  ErrNo     `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   *string   `json:"status_msg"`   // 返回状态描述
}
