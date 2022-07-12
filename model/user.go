package model

const PassWordCost = 12

// User 用户模型
type User struct {
	ID            int    `gorm:"column:id;primaryKey"`    //ID
	UserName      string `gorm:"column:user_name;unique"` //用户名
	Password      string `gorm:"column:pass_word"`        // 用户密码
	Nickname      string `gorm:"column:nick_name"`        // 用户昵称
	Status        string `gorm:"column:status"`           // 用户状态
	Avatar        string `gorm:"column:avatar;size:1000"` // 用户头像
	FollowCount   int64  `gorm:"column:follow_count"`     // 关注总数
	FollowerCount int64  `gorm:"column:follower_count"`   // 粉丝总数
}

func (User) TableName() string {
	return "user"
}
