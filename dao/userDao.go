package dao

import (
	"basictiktok/model"
	"basictiktok/util"
	"gorm.io/gorm"
	"sync"
)

type UserDaoInterface interface {
	QueryUserByName(username string) (user *model.User, err error)
	QueryUserByID(id int64) (model.User, error)
	Create(user *model.User) (err error)
	Follow(source, target int) error
	UnFollow(source, target int) error
}

type UserDao struct {
}

var (
	userDao     UserDaoInterface
	userDaoOnce sync.Once
)

func NewUserDaoInstance() UserDaoInterface {
	userDaoOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

// QueryUserByName 通过用户名查询用户信息
func (u *UserDao) QueryUserByName(username string) (user *model.User, err error) {
	user = new(model.User)
	if err = DB.Where("user_name=?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUserByID 通过用户id查询用户信息
func (u *UserDao) QueryUserByID(id int64) (model.User, error) {
	var user model.User
	result := DB.First(&user, id)
	return user, result.Error
}

// Create 创建一个新用户
func (u *UserDao) Create(user *model.User) (err error) {
	err = DB.Create(&user).Error
	return
}

// Follow 关注
func (u *UserDao) Follow(source, target int) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.User{}).Where("id=?", source).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql follow错误1\n", err)
			return err
		}
		err = DB.Model(model.User{}).Where("id=?", target).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql follow错误2\n", err)
			return err
		}
		return nil
	})
	return err
}

// UnFollow 取消关注
func (u *UserDao) UnFollow(source, target int) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.User{}).Where("id=?", source).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql unfollow错误1\n", err)
			return err
		}
		err = tx.Model(model.User{}).Where("id=?", target).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error
		if err != nil {
			util.Log().Error("mysql unfollow错误2\n", err)
			return err
		}
		return nil
	})
	return err
}
