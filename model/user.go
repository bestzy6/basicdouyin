package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//实例

// User 用户模型
type User struct {
	gorm.Model
	UserName       string //用户名
	PasswordDigest string // 用户密码
	Nickname       string // 用户昵称
	Status         string // 用户状态
	Avatar         string `gorm:"size:1000"` // 用户头像
	FollowCount    int64  // 关注总数
	FollowerCount  int64  // 粉丝总数
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 被封禁用户
	Suspend string = "suspend"
)

// CreateAUser 创建一个新用户
func CreateAUser(user *User) (err error) {
	err = DB.Create(&user).Error
	return
}

// QueryAUser 通过用户名查询用户信息
func QueryAUser(username string) (user *User, err error) {
	user = new(User)
	if err = DB.Debug().Where("username=?", username).Find(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUser 通过用户id查询用户信息
func QueryUser(id int64) (User, error) {
	var user User
	result := DB.First(&user, id)
	return user, result.Error
}

// GetUser 用ID获取用户
func GetUser(ID interface{}) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
