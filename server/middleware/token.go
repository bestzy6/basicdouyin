package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecret = "nwpu418"

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
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS256, userClaims)
	return tokenClaims.SignedString(jwtSecret)
}

func ParseToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
