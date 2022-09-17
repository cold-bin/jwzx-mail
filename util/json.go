// @author cold bin
// @date 2022/9/15

package util

import "github.com/gin-gonic/gin"

type Json struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ResOk(c *gin.Context) {
	c.JSON(200, Json{
		Code: 1000,
		Msg:  "ok",
	})
}

func ResOkWithData(c *gin.Context, data interface{}) {
	c.JSON(200, Json{
		Code: 1001,
		Msg:  "ok with data",
		Data: data,
	})
}

func ResErrMsg(c *gin.Context, msg string) {
	c.JSON(200, Json{
		Code: 1002,
		Msg:  msg,
	})
}

func ResInternalErr200(c *gin.Context) {
	c.JSON(200, Json{
		Code: 1003,
		Msg:  "server busy",
	})
}

func ResInternalErr500(c *gin.Context) {
	c.JSON(500, Json{
		Code: 1004,
		Msg:  "server busy",
	})
}
