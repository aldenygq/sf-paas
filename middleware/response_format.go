package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	Ctx *gin.Context
}

type response struct {
	Errno  int         `json:"errno"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errmsg"`
}

func (c *Context) Response(errno int, errmsg string, data interface{}) {

	c.Ctx.JSON(http.StatusOK, response{
		Errno:  errno,
		Data:   data,
		ErrMsg: errmsg,
	})
}
