package graphdb

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//测试前：数据装载、配置初始化等前置工作
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
	}
	Neo4j()
	code := m.Run()
	//测试后：释放资源等收尾工作
	Close()
	os.Exit(code)
}

func TestUser_Clear(t *testing.T) {
	ClearAll()
}

// 注意，此测试会清空数据库
func TestUser_Create(t *testing.T) {
	ClearAll()
	//CREATE (u:Users{id:100,name:'zy',follow:0,follower:0}) RETURN u
	user1 := User{
		ID:            100,
		Name:          "zy",
		FollowCount:   0,
		FollowerCount: 0,
	}
	user2 := User{
		ID:            101,
		Name:          "xj",
		FollowCount:   0,
		FollowerCount: 0,
	}
	user3 := User{
		ID:            102,
		Name:          "sb",
		FollowCount:   0,
		FollowerCount: 0,
	}
	user4 := User{
		ID:            103,
		Name:          "sb2",
		FollowCount:   0,
		FollowerCount: 0,
	}
	err := user1.Create()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	err = user2.Create()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	err = user3.Create()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	err = user4.Create()
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestUser_Follow(t *testing.T) {
	//ClearAll()
	//CREATE (u:Users{id:100,name:'zy',follow:0,follower:0}) RETURN u
	user1 := User{
		ID:   100,
		Name: "zy",
	}
	user2 := User{
		ID:   101,
		Name: "xj",
	}
	user3 := User{
		ID:   102,
		Name: "sb",
	}
	user4 := User{
		ID:   103,
		Name: "sb2",
	}
	//user1.Create()
	//user2.Create()
	err := user1.Follow(&user3)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	err = user2.Follow(&user3)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	err = user1.Follow(&user4)
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestUser_UnFollow(t *testing.T) {
	user1 := User{
		ID:   100,
		Name: "zy",
	}
	user2 := User{
		ID:   102,
		Name: "xj",
	}
	// MATCH (a:Users{id:100})-[rel:follow]->(b:Users{id:101}) DELETE rel
	// MATCH (a:Users),(b:Users) WHERE a.id=100 AND b.id=101 SET a.follow = a.follow-1, b.follower = b.follower-1 RETURN a,b
	err := user1.UnFollow(&user2)
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestUser_HasFollow(t *testing.T) {
	user1 := User{
		ID:   100,
		Name: "zy",
	}
	user2 := User{
		ID:   101,
		Name: "xj",
	}
	ok, err := user1.HasFollow(&user2)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	fmt.Println(ok)
}

func TestUser_Followers(t *testing.T) {
	user1 := User{
		ID:   100,
		Name: "zy",
	}
	user2 := User{
		ID:   101,
		Name: "xj",
	}
	// MATCH (a:Users)-[:follow]->(c:Users)<-[:follow]-(b:Users) where a.id=100 and b.id=101 return c
	followers, err := user1.Followers(&user2)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	for _, user := range followers {
		fmt.Println(user)
	}
}

func TestUser_Followees(t *testing.T) {
	user1 := User{
		ID:   100,
		Name: "zy",
	}
	user2 := User{
		ID:   101,
		Name: "xj",
	}
	// MATCH (a:Users)-[:follow]->(c:Users)<-[:follow]-(b:Users) where a.id=100 and b.id=101 return c
	followers, err := user1.Followees(&user2)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	for _, user := range followers {
		fmt.Println(user)
	}
}
