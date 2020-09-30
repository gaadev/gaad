/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package controllers

import (
	"errors"
	"gaad/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
	gin.Context
}

// 获取全部请求解析到map
func Response(c *gin.Context, code uint32, msg string, data map[string]interface{}) error {
	message := common.Result(code, msg, data)
	c.JSON(http.StatusOK, message)
	return errors.New("request meets a error")
}
