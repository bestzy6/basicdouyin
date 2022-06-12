package serializer

// LikesRequest 赞请求
type LikesRequest struct {
	Token      string `form:"token" json:"token" binding:"required"`             //用户鉴权token
	VideoId    int    `form:"video_id" json:"video_id" binding:"required"`       //视频id
	ActionType int    `form:"action_type" json:"action_type" binding:"required"` //1-点赞，2-取消点赞
}

// LikesResponse 赞响应
type LikesResponse struct {
	StatusCode ErrNo  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// LikeListRequest 点赞列表请求
type LikeListRequest struct {
	Token  string `form:"token" json:"token" binding:"required"`     //用户鉴权token
	UserId int    `form:"user_id" json:"user_id" binding:"required"` //视频id
}

// LikeListResponse 点赞列表
type LikeListResponse struct {
	StatusCode ErrNo    `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string   `json:"status_msg"`  // 返回状态描述
	VideoList  []*Video `json:"video_list"`  // 用户点赞视频列表
}

// CommentRequest 评论请求
type CommentRequest struct {
	Token       string `form:"token" json:"token" binding:"required"`               //用户鉴权token
	VideoId     int    `form:"video_id" json:"video_id" binding:"required"`         //视频id
	ActionType  int    `form:"action_type" json:"action_type" binding:"required"`   //1-发布评论，2-删除评论
	CommentText string `form:"comment_text" json:"comment_text" binding:"required"` //用户填写的评论内容，在action_type=1的时候使用
	CommentId   int    `form:"comment_id" json:"comment_id" binding:"required"`     // 要删除的评论id，在action_type=2的时候使用
}

// CommentResponse 评论操作
type CommentResponse struct {
	Comment    Comment `json:"comment"`     // 评论成功返回评论内容，不需要重新拉取整个列表
	StatusCode ErrNo   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
}

type CommentListRequest struct {
	Token   string `form:"token" json:"token" binding:"required"`       //用户鉴权token
	VideoId int    `form:"video_id" json:"video_id" binding:"required"` //视频id
}

// CommListResponse 评论列表
type CommListResponse struct {
	CommentList []*Comment `json:"comment_list"` // 评论列表
	StatusCode  ErrNo      `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string     `json:"status_msg"`   // 返回状态描述
}
