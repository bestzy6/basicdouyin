package model

import (
	"fmt"
	"testing"
	"time"
)

//func TestMain(m *testing.M) {
//	if err := Init(); err != nil {
//		os.Exit(1)
//	}
//	fmt.Println("数据库连接成功")
//	m.Run()
//}
func TestCommentPostList(t *testing.T) {
	//posts, _ := NewPostDaoInstance().QueryPostByVideoId(1)
	//for _, v := range posts {
	//	fmt.Println(v)
	//}
	post := Post{
		Id:         3,
		VideoId:    1,
		UserId:     2,
		Content:    "测试数据",
		DiggCount:  2,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := NewPostDaoInstance().CreatePost(&post)
	if err != nil {

	}
	fmt.Print("插入成功")
}
