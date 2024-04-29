package service

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/luo2pei4/base-server/internal/dao"
)

var jwtSecret = []byte("rosemary")

func LoginService(name, passwd string) (string, error) {
	userInfo, err := dao.QueryUser(name, passwd)
	if err != nil {
		return "", err
	}
	nowTime := time.Now().Local()
	expiresAt := nowTime.Add(time.Minute * 15)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": name,
		"password": passwd,
		"expires":  expiresAt.Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	userInfo.LastLoginTime = nowTime
	if err = dao.UpdateLastLoginTime(userInfo); err != nil {
		return "", err
	}
	return tokenString, nil
}
