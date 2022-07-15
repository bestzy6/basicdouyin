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
