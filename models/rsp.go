/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Rsp struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Result(code uint32, message string, data interface{}) *Rsp {

	message = GetErrorMessage(code, message)
	jsonMap := grantMap(code, message, data)

	return jsonMap
}

// 按照接口格式生成原数据数组
func grantMap(code uint32, message string, data interface{}) *Rsp {

	jsonMap := &Rsp{
		Code: code,
		Msg:  message,
		Data: data,
	}
	return jsonMap
}

func (r *Rsp) Write(c *gin.Context)  {
	c.JSON(http.StatusOK, r)
}

