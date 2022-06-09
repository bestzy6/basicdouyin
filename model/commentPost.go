package model

import (
	"basictiktok/util"
	"sync"
)

type Post struct {
	Id         int64  `gorm:"column:id"`
	VideoId    int64  `gorm:"column:video_id"`
	UserId     int64  `gorm:"column:user_id"`
	Content    string `gorm:"column:content"`
	DiggCount  int32  `gorm:"column:digg_count"`
	CreateTime string `gorm:"column:create_time"`
}

func (Post) TableName() string {
	return "post"
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

//func (*PostDao) QueryPostById(id int64) (*Post, error) {
//	var post Post
//	err := db.Where("id = ?", id).Find(&post).Error
//	if err == gorm.ErrRecordNotFound {
//		return nil, nil
//	}
//	if err != nil {
//		util.Log().Error("find video by id err:", err.Error())
//		return nil, err
//	}
//	return &post, nil
//}

// QueryPostByVideoId 根据视频Id 查询所有的评论
func (*PostDao) QueryPostByVideoId(videoId int64) ([]*Post, error) {
	var posts []*Post
	err := db.Where("video_id = ?", videoId).Find(&posts).Error
	if err != nil {
		util.Log().Error("find posts by video_id err:" + err.Error())
		return nil, err
	}
	return posts, nil
}

// CreatePost 提交评论
func (*PostDao) CreatePost(post *Post) error {
	if err := db.Create(post).Error; err != nil {
		util.Log().Error("insert post err:" + err.Error())
		return err
	}
	return nil
}
