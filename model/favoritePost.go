package model

import (
	"basictiktok/util"
	"gorm.io/gorm"
	"sync"
)

type FavoritePost struct {
	UserId    int64 `gorm:"column:user_id"`
	VideoId   int64 `gorm:"column:video_id"`
	DiggCount int32 `gorm:"column:digg_count"`
}

func (FavoritePost) TableName() string {
	return "favorite_post"
}

type FavoritePostDao struct {
}

var (
	favoritePostDao  *FavoritePostDao
	favoritePostOnce sync.Once
)

func NewFavoritePostDaoInstance() *FavoritePostDao {
	favoritePostOnce.Do(
		func() {
			favoritePostDao = &FavoritePostDao{}
		})
	return favoritePostDao
}

// QueryFavoritePostById 先用userid 查到所有的视频id,
func (*FavoritePostDao) QueryFavoritePostById(userid int64) ([]*FavoritePost, error) {
	var videoId []*FavoritePost
	err := DB.Table("douyin.favorite_post").Where("user_id = ?", userid).Find(&videoId).Error
	if err == gorm.ErrRecordNotFound {
		util.Log().Error("QueryFavoritePostById err\n", err)
		return nil, nil
	}
	if err != nil {
		util.Log().Error("find video by id err:", err.Error())
		return nil, err
	}
	return videoId, nil
}

// GetVideoIdList 获取查询得到 videoId
func (*FavoritePostDao) GetVideoIdList(videoId []*FavoritePost) []int64 {
	var videoIdList []int64
	for _, v := range videoId {
		videoIdList = append(videoIdList, v.VideoId)
	}
	return videoIdList
}

// QueryPostByUserId 根据video_Id 查询所有的点赞的视频
func (*FavoritePostDao) QueryPostByVedioId(videoLs []int64) ([]*Video, error) {
	var videos []*Video
	err := DB.Table(Video{}.TableName()).Where("id in (?)", videoLs).Find(&videos).Error // 优化地方
	if err != nil {
		util.Log().Error("find posts by video_id err:", err)
		return nil, err
	}
	return videos, nil
}

// CreateFPost 点赞
func (*FavoritePostDao) CreateFPost(fPost *FavoritePost) error {
	if err := DB.Table("douyin.favorite_post").Create(fPost).Error; err != nil {
		util.Log().Error("insert post err:" + err.Error())
		return err
	}
	return nil
}
