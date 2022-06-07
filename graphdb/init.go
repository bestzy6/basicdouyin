package graphdb

import (
	"basictiktok/util"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"os"
)

var driver neo4j.Driver

func Neo4j() {
	Driver, err := neo4j.NewDriver(os.Getenv("NEO4J_ADDR"), neo4j.BasicAuth(os.Getenv("NEO4J_USER"), os.Getenv("NEO4J_PW"), ""))
	if err != nil {
		util.Log().Panic("连接Neo4j不成功", err)
	}
	driver = Driver
	createIndex()
}

// Deprecated: 清空所有数据，慎用
func ClearAll() {
	session := newSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Println("session close failed", err)
		}
	}(session)
	_, err := session.Run("MATCH (n) OPTIONAL MATCH (n)-[r]-() DELETE n,r", map[string]interface{}{})
	if err != nil {
		util.Log().Error("clear all err:", err)
	}
}

func createIndex() {
	session := newSession()
	defer session.Close()
	_, err := session.Run("CREATE CONSTRAINT IF NOT EXISTS FOR (u:Users) REQUIRE u.id IS UNIQUE", map[string]interface{}{})
	if err != nil {
		util.Log().Error("create graph index failed:", err)
	}
}

func newSession() neo4j.Session {
	return driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
}
