package model

import (
	"gorm.io/gorm"
	"sync"
)

type Video struct {
	ID            int64  `gorm:"column:id;AUTO_INCREMENT"` // 视频唯一标识
	UserID        int64  // 视频作者id
	CommentCount  int64  // 视频的评论总数
	CoverURL      string // 视频封面地址
	FavoriteCount int64  // 视频的点赞总数
	PlayURL       string // 视频播放地址
	Title         string // 视频标题
	AddTime       int64  // 视频添加时间
}

type VideoDao struct {
}

var (
	videoDao     *VideoDao
	videoDaoOnce sync.Once
)

func NewVideoDaoInstance() *VideoDao {
	videoDaoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

func CreateAVideo(video *Video) (err error) {
	err = DB.Debug().Create(&video).Error
	return
}

func FindVideoBeforeTime(time int64) ([]*Video, error) {
	var videos []*Video
	err := DB.Debug().Where("add_time > ?", time).Limit(30).Order("add_time").Find(&videos).Error
	return videos, err
}

// QueryVideoListByUserID返回用户发布视频列表
func QueryVideoListByUserID(userid int) ([]*Video, error) {
	var videos []*Video
	err := DB.Debug().Where("user_id = ?", userid).Find(&videos).Error
	return videos, err
}

// AddComment 增加评论
func (v *VideoDao) AddComment(vid int64) error {
	err := DB.Where("id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	return err
}

// DeleteComment 删除评论
func (v *VideoDao) DeleteComment(vid int64) error {
	err := DB.Where("id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	return err
}

// AddFavorite 点赞
func (v *VideoDao) AddFavorite(vid int64) error {
	err := DB.Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	return err
}

// DeleteFavorite 取消点赞
func (v *VideoDao) DeleteFavorite(vid int64) error {
	err := DB.Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	return err
}
