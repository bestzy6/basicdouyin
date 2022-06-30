package graphdb

import (
	"basictiktok/util"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type User struct {
	ID   int    // 用户id,id
	Name string // 用户名称,name
	//FollowCount   int    // 关注总数,follow
	//FollowerCount int    // 粉丝总数,follower
	IsFollow bool // true-已关注，false-未关注
}

// Create 创建用户
func (u User) Create() error {
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
			"ID":   u.ID,
			"Name": u.Name,
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
func (u User) Favorite(v *Video) error {
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
			"UID": u.ID,
			"VID": v.ID,
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
		return nil, result.Err()
	})
	return err
}

// UnFavorite 取消点赞视频
func (u User) UnFavorite(v *Video) error {
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
			"UID": u.ID,
			"VID": v.ID,
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
		tx.Commit()
		return nil, result.Err()
	})
	return err
}

func (u User) MyFavoriteList() (map[int]*Video, error) {
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
			"ID": u.ID,
		},
	)
	if err != nil {
		return nil, err
	}
	list := make(map[int]*Video, 0)
	for result.Next() {
		record := result.Record()
		video := u.record2Video(record, "b")
		video.isFavorite = true
		list[video.ID] = video
	}
	return list, nil
}

func (u User) FavoriteList(requestor *User) (map[int]*Video, error) {
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
				"ID": u.ID,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			video := u.record2Video(record, "b")
			videos[video.ID] = video
		}
		//获取同时关注的人
		result, err = tx.Run("MATCH (a:Users)-[:favorite]->(c:Videos)<-[:favorite]-(b:Users) "+
			"where a.id=$ID and b.id=$RID "+
			"return c",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			video := u.record2Video(record, "c")
			videos[video.ID].isFavorite = true
		}
		//
		tx.Commit()
		return nil, nil
	})
	//
	return videos, err
}

// Follow 关注用户，target的ID必填
func (u User) Follow(target *User) error {
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
			"AID": u.ID,
			"BID": target.ID,
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

// UnFollow 取消关注，target的ID必填
func (u User) UnFollow(target *User) error {
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
			"AID": u.ID,
			"BID": target.ID,
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

// MyFollowers 如果请求ID与UserID相同，则返回本人的关注者
func (u User) MyFollowers() (map[int]*User, error) {
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
			"ID": u.ID,
		},
	)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		record := result.Record()
		user := u.record2User(record, "b")
		user.IsFollow = true
		users[user.ID] = user
	}
	return users, nil
}

// Followers 关注列表，requestor的ID必填
func (u User) Followers(requestor *User) (map[int]*User, error) {
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
				"ID": u.ID,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			user := u.record2User(record, "b")
			users[user.ID] = user
		}
		//获取同时关注的人
		result, err = tx.Run("MATCH (a:Users)-[:follow]->(c:Users)<-[:follow]-(b:Users) "+
			"where a.id=$ID and b.id=$RID "+
			"return c",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			user := u.record2User(record, "c")
			users[user.ID].IsFollow = true
		}
		//
		tx.Commit()
		return nil, nil
	})
	return users, err
}

// Followees 粉丝列表，requestor的ID必填
func (u User) Followees(requestor *User) (map[int]*User, error) {
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
				"ID": u.ID,
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
			user := u.record2User(record, "b")
			users[user.ID] = user
		}
		//获取是u的粉丝，且是requestor关注的人
		result, err = tx.Run("MATCH (a:Users{id:$ID})<-[:follow]-(c:Users)<-[:follow]-(b:Users{id:$RID}) "+
			"RETURN c",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
			},
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			user := u.record2User(record, "c")
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

// HasFollow 判断是否关注，，target的ID必填，true为已关注，false为未关注
func (u User) HasFollow(target *User) (bool, error) {
	//默认自己关注了自己
	if u.ID == target.ID {
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
			"AID": u.ID,
			"BID": target.ID,
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

// IsFollow 判断是否关注
func IsFollow(src, target int) (bool, error) {
	//同一个人时，返回真
	if src == target {
		return true, nil
	}
	//
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
			"AID": src,
			"BID": target,
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
func (u User) record2User(record *neo4j.Record, key string) *User {
	get, _ := record.Get(key)
	node := get.(dbtype.Node)
	return &User{
		ID:   int(node.Props["id"].(int64)),
		Name: node.Props["name"].(string),
	}
}

//内部函数，记录转为video
func (u User) record2Video(record *neo4j.Record, key string) *Video {
	get, _ := record.Get(key)
	node := get.(dbtype.Node)
	return &Video{
		ID:            int(node.Props["id"].(int64)),
		UserID:        int(node.Props["userid"].(int64)),
		CoverURL:      node.Props["coverUrl"].(string),
		FavoriteCount: int(node.Props["favoriteCount"].(int64)),
	}
}
