package model

import (
	"basictiktok/util"
	"gorm.io/gorm"
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
func (v *VideoClDao) QueryByVideoId(videoId int64) (*VideoCL, error) {
	var video VideoCL
	err := DB.Table("douyin.videoCL").Where("video_id = ?", videoId).Find(&video).Error
	if err != nil {
		util.Log().Error("find posts by video_id err:" + err.Error())
		return nil, err
	}
	return &video, nil
}

// DeFavorite 点赞-1
func (v *VideoClDao) DeFavorite(vid int64) error {
	err := DB.Table("douyin.videoCL").Where("video_id=?", vid).UpdateColumn("favorite_count", gorm.Expr("follower_count - ?", 1)).Error
	return err
}

// AddFavorite 点赞+1
func (v *VideoClDao) AddFavorite(vid int64) error {
	err := DB.Table("douyin.videoCL").Where("video_id=?", vid).UpdateColumn("favorite_count", gorm.Expr("follower_count + ?", 1)).Error
	return err
}

// AddComment 评论数+1
func (v *VideoClDao) AddComment(vid int64) error {
	err := DB.Table("douyin.videoCL").Where("video_id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	return err
}
