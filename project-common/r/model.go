package r

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BusinessCode int

type Result struct {
	Code BusinessCode `json:"code"`
	Msg  string       `json:"msg"`
	Data any          `json:"data"`
}

func Success(ctx *gin.Context, data ...any) {
	r := &Result{
		Code: 200,
		Msg:  "success",
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	ctx.JSON(http.StatusOK, r)
}

func Fail(ctx *gin.Context, code int, msg string) {
	ctx.JSON(http.StatusOK, &Result{
		Code: BusinessCode(code),
		Msg:  msg},
	)
}
