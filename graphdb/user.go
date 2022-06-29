package graphdb

import (
	"basictiktok/util"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type User struct {
	ID            int    // 用户id,id
	Name          string // 用户名称,name
	FollowCount   int    // 关注总数,follow
	FollowerCount int    // 粉丝总数,follower
	IsFollow      bool   // true-已关注，false-未关注
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
	result, err := session.Run("CREATE (u:Users{id:$ID,name:$Name,follow:$FollowCount,follower:$FollowerCount}) "+
		"RETURN u.id,u.name",
		map[string]interface{}{
			"ID":            u.ID,
			"Name":          u.Name,
			"FollowCount":   u.FollowCount,
			"FollowerCount": u.FollowerCount,
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
		//更新双方结点的属性值
		result, err = tx.Run(
			"MATCH (a:Users),(b:Users) "+
				"WHERE a.id=$AID AND b.id=$BID "+
				"SET a.follow = a.follow+1, b.follower = b.follower+1 "+
				"RETURN a,b", hash)
		if err != nil {
			return nil, err
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
		//更新双方结点的属性值
		result, err = tx.Run("MATCH (a:Users),(b:Users) WHERE a.id=$AID AND b.id=$BID "+
			"SET a.follow = a.follow-1, b.follower = b.follower-1 "+
			"RETURN a,b", hash)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		record, err = result.Single()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		a, _ := record.Get("a")
		b, _ := record.Get("b")
		nodeA, nodeB := a.(dbtype.Node), b.(dbtype.Node)
		util.Log().Info("取消关注成功！用户[ID:%v Name:%v]的关注数-1，用户[ID:%v Name:%v]的粉丝数-1\n",
			nodeA.Props["id"],
			nodeA.Props["name"],
			nodeB.Props["id"],
			nodeB.Props["name"])
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
	//users := make([]*User, 0, u.FollowerCount)
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
		fmt.Println(requestor.ID, u.ID)
		fmt.Println("hello")
		for result.Next() {
			record := result.Record()
			user := u.record2User(record, "c")
			fmt.Println(user.ID)
			users[user.ID].IsFollow = true
			//users = append(users, user)
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

//记录转为user
func (u User) record2User(record *neo4j.Record, key string) *User {
	get, _ := record.Get(key)
	node := get.(dbtype.Node)
	return &User{
		ID:            int(node.Props["id"].(int64)),
		Name:          node.Props["name"].(string),
		FollowCount:   int(node.Props["follow"].(int64)),
		FollowerCount: int(node.Props["follower"].(int64)),
	}
}

// IsFollow 判断是否关注
func IsFollow(src, target int) (bool, error) {
	//同一个人时，返回真
	if src == target {
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
