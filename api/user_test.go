package api

import (
	"basictiktok/serializer"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

//func Register(c *gin.Context){
//	var req serializer.RegisterRequest
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(http.StatusOK, serializer.RegisterResponse{
//			StatusCode: serializer.ParamInvalid,
//			StatusMsg:  "请求参数错误",
//		})
//		return
//	}
//	resp := service.RegisterService(&req)
//	c.JSON(http.StatusOK,resp)
//}
func SetupServer() *gin.Engine {
	router := gin.Default() // 这需要写到init中，启动gin框架
	router.POST("/register", Register)
	return router //把启动的engine 对象传⼊到test框架中
}

// Register单元测试
func TestRegister(t *testing.T) {

	uri := "/register"
	param := url.Values{
		"nickname":        {"hrewq"},
		"user_name":       {"hello"},
		"password_digest": {"1321432"},
	}
	router := SetupServer()
	body := PostForm(uri, param, router)
	resp := &serializer.RegisterResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		t.Errorf("解析响应出错,err:%v\n", err)
	}
	t.Log(resp.StatusMsg)
	if resp.StatusCode != serializer.OK {
		t.Errorf("响应数据不符，errmsg:%v", resp.StatusMsg)
	}
}

func PostForm(uri string, param url.Values, router *gin.Engine) []byte {
	// 构造post请求
	req := httptest.NewRequest("POST", uri, strings.NewReader(param.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调⽤相应handler接⼝
	router.ServeHTTP(w, req)
	// 提取响应
	result := w.Result()
	defer result.Body.Close()
	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}
