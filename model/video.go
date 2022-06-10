package model

import (
	"time"
)

type Video struct {
	ID            int64     `gorm:"column:id;AUTO_INCREMENT"` // 视频唯一标识
	UserID        int64     // 视频作者id
	CommentCount  int64     // 视频的评论总数
	CoverURL      string    // 视频封面地址
	FavoriteCount int64     // 视频的点赞总数
	PlayURL       string    // 视频播放地址
	Title         string    // 视频标题
	AddTime       time.Time // 视频添加时间
}

func CreateAVideo(video *Video) (err error) {
	err = DB.Debug().Create(&video).Error
	return
}
