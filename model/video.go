package model

type Video struct {
	ID            int64  `gorm:"column:id;primaryKey"`  // 视频唯一标识
	UserID        int64  `gorm:"column:user_id"`        // 视频作者id
	CommentCount  int64  `gorm:"column:comment_count"`  // 视频的评论总数
	CoverURL      string `gorm:"column:cover_url"`      // 视频封面地址
	FavoriteCount int64  `gorm:"column:favorite_count"` // 视频的点赞总数
	PlayURL       string `gorm:"column:play_url"`       // 视频播放地址
	Title         string `gorm:"column:title"`          // 视频标题
	AddTime       int64  `gorm:"column:add_time"`       // 视频添加时间
}

func (Video) TableName() string {
	return "video"
}
