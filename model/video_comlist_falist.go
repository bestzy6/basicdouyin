package model

import (
	"basictiktok/util"
	"sync"
)

type VideoCL struct {
	VideoId       int64 `gorm:"column:video_id"`
	CommentCount  int64 `gorm:"column:comment_count"`  // 视频的评论总数
	FavoriteCount int64 `gorm:"column:favorite_count"` // 视频的点赞总数
}

func (VideoCL) TableName() string {
	return "videoCL"
}

type VideoClDao struct {
}

var (
	videoClDao     *VideoClDao
	videoClDaoOnce sync.Once
)

func NewVideoClDaoInstance() *VideoClDao {
	videoClDaoOnce.Do(
		func() {
			videoClDao = &VideoClDao{}
		})
	return videoClDao
}
func (v VideoClDao) QueryByVideoId(videoId int64) (*VideoCL, error) {
	var video VideoCL
	err := DB.Table("douyin.videoCL").Where("video_id = ?", videoId).Find(&video).Error
	if err != nil {
		util.Log().Error("find posts by video_id err:" + err.Error())
		return nil, err
	}
	return &video, nil
}
