package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecret = []byte("nwpu418")

type UserClaims struct {
	*jwt.RegisteredClaims
	UserID int
}

func CreateToken(userID int) (string, error) {
	expiresTime := jwt.NumericDate{Time: time.Now().Add(time.Hour * 10)}
	userClaims := &UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: &expiresTime,
		},
		UserID: userID,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return tokenClaims.SignedString(jwtSecret)
}

// ParseToken 解析token，如果过期会返回错误
func ParseToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		Log().Error("解析token出错\n", err)
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.New("解析token出错")
}
