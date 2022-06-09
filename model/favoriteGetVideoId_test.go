package model

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := Init(); err != nil {
		os.Exit(1)
	}
	fmt.Println("数据库连接成功")
	m.Run()
}

func TestFavoriteGetVideoIdByUserId(t *testing.T) {
	post := FavoritePost{
		VideoId:   3,
		UserId:    1,
		DiggCount: 2,
	}
	err := NewFavoritePostDaoInstance().CreateFPost(&post)
	if err != nil {
	}
	fmt.Println("插入成功")

	//v := NewFavoritePostDaoInstance() // 查询user_id == 1 用户点赞的视频id
	//b, _ := v.QueryFavoritePostById(1)
	//c := v.GetVideoIdList(b)
	//d, _ := v.QueryPostByUserId(c)
	//for _, l := range b {
	//	fmt.Println(l)
	//}

}
