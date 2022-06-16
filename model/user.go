package model

import (
	"basictiktok/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const PassWordCost = 12

// User 用户模型
type User struct {
	ID            int    `gorm:"column:id;AUTO_INCREMENT"`
	UserName      string `gorm:"unique"` //用户名
	Password      string // 用户密码
	Nickname      string // 用户昵称
	Status        string // 用户状态
	Avatar        string `gorm:"size:1000"` // 用户头像
	FollowCount   int64  // 关注总数
	FollowerCount int64  // 粉丝总数
}

// QueryUserByName 通过用户名查询用户信息
func QueryUserByName(username string) (user *User, err error) {
	user = new(User)
	if err = DB.Where("user_name=?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUserByID 通过用户id查询用户信息
func QueryUserByID(id int64) (User, error) {
	var user User
	result := DB.First(&user, id)
	return user, result.Error
}

// Create 创建一个新用户
func (user *User) Create() (err error) {
	err = DB.Create(&user).Error
	return
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// Follow 关注
func (user *User) Follow(toUser *User) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(user).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql follow错误1\n", err)
			return err
		}
		err = DB.Model(toUser).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql follow错误2\n", err)
			return err
		}
		return nil
	})
	return err
}

// UnFollow 取消关注
func (user *User) UnFollow(toUser *User) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(user).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql unfollow错误1\n", err)
			return err
		}
		err = tx.Model(toUser).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql unfollow错误2\n", err)
			return err
		}
		return nil
	})
	return err
}
