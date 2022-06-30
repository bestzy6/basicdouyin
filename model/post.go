package model

import (
	"basictiktok/util"
	"gorm.io/gorm"
	"sync"
)

type Post struct {
	Id      int64  `gorm:"column:id;primaryKey;AUTO_INCREMENT""`
	VideoId int64  `gorm:"column:video_id"`
	UserId  int64  `gorm:"column:user_id"`
	Content string `gorm:"column:content"`
	//DiggCount  int32  `gorm:"column:digg_count"`
	CreateTime string `gorm:"column:create_time"`
}

func (Post) TableName() string {
	return "post"
}

func (p *Post) Creat() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Create(p).Error
		if err != nil {
			return err
		}
		err = tx.Model(&Video{}).Where("id=?", p.VideoId).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		return err
	})
}

func (p *Post) Delete() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Delete(p).Error
		if err != nil {
			return err
		}
		err = tx.Model(&Video{}).Where("id=?", p.VideoId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		return err
	})
}

type PostDao struct {
}

var postDao *PostDao
var postOnce sync.Once

func NewPostDaoInstance() *PostDao {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}

// QueryPostByVideoId 根据视频Id 查询所有的评论
func (*PostDao) QueryPostByVideoId(videoId int64) ([]*Post, error) {
	var posts []*Post
	err := DB.Table(Post{}.TableName()).Where("video_id = ?", videoId).Find(&posts).Error
	if err != nil {
		util.Log().Error("find posts by video_id err:" + err.Error())
		return nil, err
	}
	return posts, nil
}

// CreatePost 提交评论
//func (*PostDao) CreatePost(post *Post) error {
//	if err := DB.Table(Post{}.TableName()).Create(post).Error; err != nil {
//		util.Log().Error("insert post err:" + err.Error())
//		return err
//	}
//	return nil
//}
