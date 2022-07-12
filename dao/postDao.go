package dao

import (
	"basictiktok/model"
	"basictiktok/util"
	"gorm.io/gorm"
	"sync"
)

type PostDaoInterface interface {
	QueryPostByVideoId(videoId int64) ([]*model.Post, error)
	Creat(p *model.Post) error
	Delete(videoId, commentId int) error
}

type PostDao struct {
}

var (
	postDao  PostDaoInterface
	postOnce sync.Once
)

func NewPostDaoInstance() PostDaoInterface {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}

// QueryPostByVideoId 根据视频Id 查询所有的评论
func (*PostDao) QueryPostByVideoId(videoId int64) ([]*model.Post, error) {
	var posts []*model.Post
	err := DB.Table(model.Post{}.TableName()).Where("video_id = ?", videoId).Find(&posts).Error
	if err != nil {
		util.Log().Error("find posts by video_id err:" + err.Error())
		return nil, err
	}
	return posts, nil
}

func (*PostDao) Creat(p *model.Post) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Create(p).Error
		if err != nil {
			return err
		}
		err = tx.Model(&model.Video{}).Where("id=?", p.VideoId).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		return err
	})
}

func (*PostDao) Delete(videoId, commentId int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Delete(model.Post{}, commentId).Error
		if err != nil {
			return err
		}
		err = tx.Model(&model.Video{}).Where("id=?", videoId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		return err
	})
}
