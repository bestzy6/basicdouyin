package graphdb

import (
	"basictiktok/util"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"sync"
)

type User struct {
	ID       int    // 用户id,id
	Name     string // 用户名称,name
	IsFollow bool   // true-已关注，false-未关注
}

type UserGraphDao struct {
}

var (
	userGraphDao     *UserGraphDao
	userGraphDaoOnce sync.Once
)

func NewUserGraphDao() *UserGraphDao {
	userGraphDaoOnce.Do(func() {
		userGraphDao = new(UserGraphDao)
	})
	return userGraphDao
}

func (u *UserGraphDao) Create(user *User) error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	result, err := session.Run("CREATE (u:Users{id:$ID,name:$Name}) "+
		"RETURN u.id,u.name",
		map[string]interface{}{
			"ID":   user.ID,
			"Name": user.Name,
		},
	)
	if err != nil {
		return err
	}
	record, err := result.Single()
	if err != nil {
		return err
	}
	id, _ := record.Get("u.id")
	name, _ := record.Get("u.name")
	util.Log().Info("创建用户图结点[ID:%v Name:%v]成功!\n", id, name)
	return result.Err()
}

// Favorite 点赞视频
func (u *UserGraphDao) Favorite(userId, videoId int) error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		hash := map[string]interface{}{
			"UID": userId,
			"VID": videoId,
		}
		//查询是否存在关系
		result, err := tx.Run(
			"MATCH (a:Users{id:$UID})-[rel:favorite]->(b:Videos{id:$VID}) "+
				"RETURN COUNT(rel)", hash)
		if err != nil {
			return nil, err
		}
		//如果已经存在关注关系，则返回错误
		record, err := result.Single()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		n, b := record.Get("COUNT(rel)")
		if b {
			fmt.Println(n)
		}
		if n, _ := record.Get("COUNT(rel)"); n.(int64) > 0 {
			tx.Rollback()
			return nil, errors.New("点赞失败，已点赞！")
		}
		//建立关注关系
		result, err = tx.Run(
			"MATCH (a:Users{id:$UID}),(b:Videos{id:$VID}) "+
				"CREATE (a)-[r1:favorite]->(b) "+
				"RETURN r1", hash)
		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}
		//视频获赞数+1
		result, err = tx.Run(
			"MATCH (b:Videos{id:$VID}) SET b.favoriteCount = b.favoriteCount+1", hash)
		if err != nil {
			return nil, err
		}
		tx.Commit()
		return nil, result.Err()
	})
	return err
}

// UnFavorite 取消点赞视频
func (u *UserGraphDao) UnFavorite(userId, videoId int) error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		hash := map[string]interface{}{
			"UID": userId,
			"VID": videoId,
		}
		//取消关注关系
		result, err := tx.Run(
			"MATCH (a:Users{id:$UID})-[rel:favorite]->(b:Videos{id:$VID}) "+
				"DELETE rel "+
				"RETURN COUNT(rel)", hash)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		//如果不存在关注关系，则返回错误
		record, err := result.Single()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if n, _ := record.Get("COUNT(rel)"); n.(int64) <= 0 {
			tx.Rollback()
			return nil, errors.New("取消点赞失败，没有点赞！")
		}
		//视频获赞数-1
		result, err = tx.Run(
			"MATCH (b:Videos{id:$VID}) SET b.favoriteCount = b.favoriteCount-1", hash)
		if err != nil {
			return nil, err
		}
		tx.Commit()
		return nil, result.Err()
	})
	return err
}

func (u *UserGraphDao) MyFavoriteList(userid int) (map[int]*Video, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	result, err := session.Run("MATCH (a:Users{id:$ID})-[:favorite]->(b:Videos) "+
		"RETURN b",
		map[string]interface{}{
			"ID": userid,
		},
	)
	if err != nil {
		return nil, err
	}
	list := make(map[int]*Video, 0)
	for result.Next() {
		record := result.Record()
		video := record2Video(record, "b")
		video.isFavorite = true
		list[video.ID] = video
	}
	return list, nil
}

func (u *UserGraphDao) FavoriteList(reqUserId, userId int) (map[int]*Video, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	videos := make(map[int]*Video)
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		//获取u关注的人
		result, err := tx.Run("MATCH (a:Users{id:$ID})-[:favorite]->(b:Videos) "+
			"RETURN b",
			map[string]interface{}{
				"ID": userId,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			video := record2Video(record, "b")
			videos[video.ID] = video
		}
		//获取同时关注的人
		result, err = tx.Run("MATCH (a:Users)-[:favorite]->(c:Videos)<-[:favorite]-(b:Users) "+
			"where a.id=$ID and b.id=$RID "+
			"return c",
			map[string]interface{}{
				"ID":  userId,
				"RID": reqUserId,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			video := record2Video(record, "c")
			videos[video.ID].isFavorite = true
		}
		//
		tx.Commit()
		return nil, nil
	})
	//
	return videos, err
}

func (u *UserGraphDao) IsFavorite(userId, videoId int) bool {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	result, err := session.Run("MATCH (a:Users{id:$UID})-[rel:favorite]->(b:Videos{id:$VID}) "+
		"RETURN COUNT(rel)",
		map[string]interface{}{
			"UID": userId,
			"VID": videoId,
		})
	if err != nil {
		util.Log().Error("IsFavorite执行出错！", err)
		return false
	}
	record, err := result.Single()
	if err != nil {
		util.Log().Error("IsFavorite执行出错！", err)
		return false
	}
	n, _ := record.Get("COUNT(rel)")
	return n.(int64) > 0
}

func (u *UserGraphDao) Follow(sourceId, targetId int) error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		hash := map[string]interface{}{
			"AID": sourceId,
			"BID": targetId,
		}
		//查询是否存在关系
		result, err := tx.Run(
			"MATCH (a:Users{id:$AID})-[rel:follow]->(b:Users{id:$BID}) "+
				"RETURN COUNT(rel)", hash)
		if err != nil {
			return nil, err
		}
		//如果已经存在关注关系，则返回错误
		record, err := result.Single()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if n, _ := record.Get("COUNT(rel)"); n.(int64) > 0 {
			tx.Rollback()
			return nil, errors.New("关注失败，已关注！")
		}
		//建立关注关系
		result, err = tx.Run(
			"MATCH (a:Users{id:$AID}),(b:Users{id:$BID}) "+
				"CREATE (a)-[r1:follow]->(b) "+
				"RETURN r1", hash)
		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}
		return nil, result.Err()
	})
	return err
}

func (u *UserGraphDao) UnFollow(sourceId, targetId int) error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		hash := map[string]interface{}{
			"AID": sourceId,
			"BID": targetId,
		}
		//取消关注关系
		result, err := tx.Run(
			"MATCH (a:Users{id:$AID})-[rel:follow]->(b:Users{id:$BID}) "+
				"DELETE rel "+
				"RETURN COUNT(rel)", hash)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		//如果不存在关注关系，则返回错误
		record, err := result.Single()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if n, _ := record.Get("COUNT(rel)"); n.(int64) <= 0 {
			tx.Rollback()
			return nil, errors.New("取消关注失败，没有关注！")
		}
		tx.Commit()
		return nil, result.Err()
	})
	return err
}

func (u *UserGraphDao) MyFollowers(userid int) (map[int]*User, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	users := make(map[int]*User)
	result, err := session.Run("MATCH (a:Users{id:$ID})-[:follow]->(b:Users) "+
		"RETURN b",
		map[string]interface{}{
			"ID": userid,
		},
	)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		record := result.Record()
		user := record2User(record, "b")
		user.IsFollow = true
		users[user.ID] = user
	}
	return users, nil
}

func (u *UserGraphDao) Followers(userId, reqUserid int) (map[int]*User, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	users := make(map[int]*User)
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		//获取u关注的人
		result, err := tx.Run("MATCH (a:Users{id:$ID})-[:follow]->(b:Users) "+
			"RETURN b",
			map[string]interface{}{
				"ID": userId,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			user := record2User(record, "b")
			users[user.ID] = user
		}
		//获取同时关注的人
		result, err = tx.Run("MATCH (a:Users)-[:follow]->(c:Users)<-[:follow]-(b:Users) "+
			"where a.id=$ID and b.id=$RID "+
			"return c",
			map[string]interface{}{
				"ID":  userId,
				"RID": reqUserid,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			user := record2User(record, "c")
			users[user.ID].IsFollow = true
		}
		//
		tx.Commit()
		return nil, nil
	})
	return users, err
}

func (u *UserGraphDao) Followees(userId, reqUserId int) (map[int]*User, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	users := make(map[int]*User)
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		//获取关注u的人
		result, err := tx.Run("MATCH (a:Users{id:$ID})<-[:follow]-(b:Users) "+
			"RETURN b",
			map[string]interface{}{
				"ID": userId,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if result.Err() != nil {
			tx.Rollback()
			return nil, result.Err()
		}
		for result.Next() {
			record := result.Record()
			user := record2User(record, "b")
			users[user.ID] = user
		}
		//获取是u的粉丝，且是requestor关注的人
		result, err = tx.Run("MATCH (a:Users{id:$ID})<-[:follow]-(c:Users)<-[:follow]-(b:Users{id:$RID}) "+
			"RETURN c",
			map[string]interface{}{
				"ID":  userId,
				"RID": reqUserId,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			user := record2User(record, "c")
			users[user.ID].IsFollow = true
		}
		if result.Err() != nil {
			tx.Rollback()
			return nil, result.Err()
		}
		tx.Commit()
		return nil, nil
	})
	return users, err
}

func (u *UserGraphDao) HasFollow(sourceId, targetId int) (bool, error) {
	//默认自己关注了自己
	if sourceId == targetId {
		return true, nil
	}
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	result, err := session.Run(
		"MATCH (a:Users{id:$AID})-[rel:follow]->(b:Users{id:$BID}) "+
			"RETURN COUNT(rel)",
		map[string]interface{}{
			"AID": sourceId,
			"BID": targetId,
		})
	if err != nil {
		util.Log().Error("HasFollow执行出错！", err)
		return false, err
	}
	record, err := result.Single()
	if err != nil {
		util.Log().Error("HasFollow执行出错！", err)
		return false, err
	}
	n, _ := record.Get("COUNT(rel)")
	return n.(int64) > 0, nil
}

//内部函数，记录转为user
func record2User(record *neo4j.Record, key string) *User {
	get, _ := record.Get(key)
	node := get.(dbtype.Node)
	return &User{
		ID:   int(node.Props["id"].(int64)),
		Name: node.Props["name"].(string),
	}
}

//内部函数，记录转为video
func record2Video(record *neo4j.Record, key string) *Video {
	get, _ := record.Get(key)
	node := get.(dbtype.Node)
	return &Video{
		ID:            int(node.Props["id"].(int64)),
		UserID:        int(node.Props["userid"].(int64)),
		CoverURL:      node.Props["coverUrl"].(string),
		FavoriteCount: int(node.Props["favoriteCount"].(int64)),
	}
}
