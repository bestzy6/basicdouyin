package graphdb

import (
	"basictiktok/util"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"sync"
)

type Video struct {
	ID            int
	UserID        int
	CoverURL      string
	FavoriteCount int
	isFavorite    bool
}

type VideoGraphDao struct {
}

var (
	videoGraphDao *VideoGraphDao
	once          sync.Once
)

func NewVideoGraphDao() *VideoGraphDao {
	once.Do(func() {
		videoGraphDao = new(VideoGraphDao)
	})
	return videoGraphDao
}

func (v *VideoGraphDao) Create(video *Video) error {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			util.Log().Error("close session err", err)
		}
	}(session)
	//
	result, err := session.Run("CREATE (v:Videos{id:$ID,userid:$UserID,coverUrl:$CoverURL,favoriteCount:$FavoriteCount}) "+
		"RETURN v.id",
		map[string]interface{}{
			"ID":            video.ID,
			"UserID":        video.UserID,
			"CoverURL":      video.CoverURL,
			"FavoriteCount": video.FavoriteCount,
		},
	)
	if err != nil {
		return err
	}
	record, err := result.Single()
	if err != nil {
		return err
	}
	id, _ := record.Get("v.id")
	util.Log().Info("创建用户图结点[ID:%v Name:%v]成功!\n", id)
	return result.Err()
}
