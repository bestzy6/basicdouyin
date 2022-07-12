package dao

import (
	"basictiktok/model"
	"gorm.io/gorm"
	"sync"
)

type VideoDaoInterface interface {
	FindVideoBeforeTime(time int64) ([]*model.Video, error)
	QueryVideoListByUserID(userid int) ([]*model.Video, error)
	Create(video *model.Video) error
	AddComment(vid int64) error
	DeleteComment(vid int64) error
	AddFavorite(vid int64) error
	DeleteFavorite(vid int64) error
}

type VideoDao struct {
}

var (
	videoDao     VideoDaoInterface
	videoDaoOnce sync.Once
)

func NewVideoDaoInstance() VideoDaoInterface {
	videoDaoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (v *VideoDao) FindVideoBeforeTime(time int64) ([]*model.Video, error) {
	var videos []*model.Video
	err := DB.Debug().Where("add_time > ?", time).Limit(30).Order("add_time").Find(&videos).Error
	return videos, err
}

// QueryVideoListByUserID返回用户发布视频列表
func (v *VideoDao) QueryVideoListByUserID(userid int) ([]*model.Video, error) {
	var videos []*model.Video
	err := DB.Where("user_id = ?", userid).Find(&videos).Error
	return videos, err
}

func (v *VideoDao) Create(video *model.Video) error {
	return DB.Create(video).Error
}

// AddComment 增加评论
func (v *VideoDao) AddComment(vid int64) error {
	err := DB.Model(&model.Video{}).Where("id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	return err
}

// DeleteComment 删除评论
func (v *VideoDao) DeleteComment(vid int64) error {
	err := DB.Model(&model.Video{}).Where("id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	return err
}

// AddFavorite 点赞
func (v *VideoDao) AddFavorite(vid int64) error {
	err := DB.Model(&model.Video{}).Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	return err
}

// DeleteFavorite 取消点赞
func (v *VideoDao) DeleteFavorite(vid int64) error {
	err := DB.Model(&model.Video{}).Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	return err
}
