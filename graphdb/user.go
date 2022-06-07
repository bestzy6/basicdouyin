package graphdb

import (
	"basictiktok/util"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type User struct {
	ID            int    // 用户id,id
	Name          string // 用户名称,name
	FollowCount   int    // 关注总数,follow
	FollowerCount int    // 粉丝总数,follower
	IsFollow      bool   // true-已关注，false-未关注
}

func (u User) Create() error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	result, err := session.Run("CREATE (u:Users{id:$ID,name:$Name,follow:$FollowCount,follower:$FollowerCount}) RETURN u",
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
	return result.Err()
}

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
		//建立关注关系
		result, err := tx.Run("MATCH (a:Users),(b:Users) WHERE a.id=$AID AND b.id=$BID "+
			"CREATE (a)-[r1:follow]->(b) "+
			"RETURN r1", hash)
		if err != nil {
			return nil, err
		}
		if result.Err() != nil {
			return nil, result.Err()
		}
		//更新双方结点的属性值
		result, err = tx.Run("MATCH (a:Users),(b:Users) WHERE a.id=$AID AND b.id=$BID "+
			"SET a.follow = a.follow+1, b.follower = b.follower+1"+
			"RETURN a,b", hash)
		if err != nil {
			return nil, err
		}
		return nil, result.Err()
	})
	return err
}

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
		result, err := tx.Run("MATCH (a:Users{id=$AID})-[rel:follow]->(b:Users{id=$BID}) DELETE rel", hash)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if result.Err() != nil {
			tx.Rollback()
			return nil, result.Err()
		}
		//更新双方结点的属性值
		result, err = tx.Run("MATCH (a:Users),(b:Users) WHERE a.id=$AID AND b.id=$BID "+
			"SET a.follow = a.follow-1, b.follower = b.follower-1"+
			"RETURN a,b", hash)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return nil, result.Err()
	})
	return err
}

// Followers 关注列表
func (u User) Followers(requestor *User) ([]*User, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	users := make([]*User, 0, u.FollowCount)
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		//获取同时关注的人
		result, err := tx.Run("MATCH (a:Users{id:$ID})-[r:follow]->(b:Users),(c:Users{id:$RID})-[r2:follow]->(d:Users) "+
			"WHERE b.id = d.id"+
			"RETURN b",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
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
			user := u.record2User(record, true)
			users = append(users, user)
		}
		//获取只有u关注的人
		result, err = tx.Run("MATCH (a:Users{id:$ID})-[r:follow]->(b:Users),(c:Users{id:$RID})-[r2:follow]->(d:Users) "+
			"WHERE b.id <> d.id"+
			"RETURN b",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
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
			user := u.record2User(record, false)
			users = append(users, user)
		}
		tx.Commit()
		return nil, nil
	})
	return users, err
}

// Followees 粉丝列表
func (u User) Followees(requestor *User) ([]*User, error) {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	users := make([]*User, 0, u.FollowerCount)
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		//获取是u的粉丝，且是requestor关注的人
		result, err := tx.Run("MATCH (a:Users{id:$ID})<-[r:follow]-(b:Users),(c:Users{id:$RID})-[r2:follow]->(d:Users) "+
			"WHERE b.id = d.id"+
			"RETURN b",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
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
			user := u.record2User(record, true)
			users = append(users, user)
		}
		//获取仅关注u的人
		result, err = tx.Run("MATCH (a:Users{id:$ID})-[r:follow]->(b:Users),(c:Users{id:$RID})-[r2:follow]->(d:Users) "+
			"WHERE b.id <> d.id"+
			"RETURN b",
			map[string]interface{}{
				"ID":  u.ID,
				"RID": requestor.ID,
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
			user := u.record2User(record, false)
			users = append(users, user)
		}
		tx.Commit()
		return nil, nil
	})
	return users, err
}

//记录转为user
func (u User) record2User(record *neo4j.Record, isFollow bool) *User {
	id, _ := record.Get("id")
	name, _ := record.Get("name")
	follow, _ := record.Get("follow")
	follower, _ := record.Get("follower")
	return &User{
		ID:            id.(int),
		Name:          name.(string),
		FollowCount:   follow.(int),
		FollowerCount: follower.(int),
		IsFollow:      isFollow,
	}
}
