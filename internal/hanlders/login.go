package hanlders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luo2pei4/base-server/internal/service"
)

type Login struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

func LoginHandler(ctx *gin.Context) {
	login := &Login{}
	if err := ctx.BindJSON(login); err != nil {
		ctx.String(http.StatusBadRequest, "parse name&passwd failed, %s", err.Error())
		return
	}
	token, err := service.LoginService(login.Name, login.Passwd)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "login failed, %s", err.Error())
		return
	}
	ctx.Header("token", token)
}
