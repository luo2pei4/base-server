package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luo2pei4/base-server/internal/service"
)

func CheckAuth(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	if len(tokenString) == 0 {
		ctx.String(http.StatusBadRequest, "illegal operation")
		return
	}
	_, err := jwt.ParseWithClaims(tokenString, &service.BaseClaims{}, func(t *jwt.Token) (interface{}, error) {
		return service.JwtSecret, nil
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Next()
}
