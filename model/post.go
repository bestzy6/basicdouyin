package model

type Post struct {
	Id         int64  `gorm:"column:id;primaryKey"`
	VideoId    int64  `gorm:"column:video_id"`
	UserId     int64  `gorm:"column:user_id"`
	Content    string `gorm:"column:content"`
	CreateTime string `gorm:"column:create_time"`
}

func (Post) TableName() string {
	return "post"
}
