package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserID        int64     // 视频作者信息
	CommentCount  int64     // 视频的评论总数
	CoverURL      string    // 视频封面地址
	FavoriteCount int64     // 视频的点赞总数
	ID            int64     // 视频唯一标识
	IsFavorite    bool      // true-已点赞，false-未点赞
	PlayURL       string    // 视频播放地址
	Title         string    // 视频标题
}

