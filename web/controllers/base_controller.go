/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package controllers

import (
	"gaad/models"
	"github.com/gin-gonic/gin"
)


type BaseController struct {
	gin.Context
}

// 获取全部请求解析到map
func Text()  {

}

// 获取全部请求解析到map
func Response(code uint32, msg string, data interface{}) *models.Rsp  {
	return models.Result(code, msg, data)
}
