package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luo2pei4/base-server/internal/dao"
	"github.com/luo2pei4/base-server/internal/utils"
)

type BaseClaims struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

var (
	JwtSecret []byte
	dur       time.Duration
)

func init() {
	JwtSecret = utils.GetSecretKey(64)
	dur = time.Minute * 15
}

func LoginService(name, passwd string) (string, error) {
	userInfo, err := dao.QueryUser(name, passwd)
	if err != nil {
		return "", err
	}
	nowTime := time.Now().Local()
	expiresAt := nowTime.Add(dur)
	claims := BaseClaims{
		name,
		passwd,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}
	userInfo.LastLoginTime = nowTime
	if err = dao.UpdateLastLoginTime(userInfo); err != nil {
		return "", err
	}
	return tokenString, nil
}
