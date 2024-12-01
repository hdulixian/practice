package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func Success(c *gin.Context, code int, msg string, data any) {
	c.JSON(http.StatusOK, Response{Code: code, Msg: msg, Data: data})
}

func Fail(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, Response{Code: code, Msg: msg})
}
