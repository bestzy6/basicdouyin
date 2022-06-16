package util

import (
	"bytes"
	"crypto/md5"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	IMG   = "./static/img"
	VEDIO = "./static/video"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// PasswordWithMD5 返回通过MD5加密之后的密码
func PasswordWithMD5(password string) string {
	data := []byte(password) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

// ReadFrameAsJpeg 从视频中截取指定帧
func ReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		fmt.Println("error")
	}
	return buf
}
