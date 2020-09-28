/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:20
 */

package routers

import (
	"fmt"
	"gaad/controllers/project"
	"gaad/initialize"
	"github.com/gin-gonic/gin"
	"net/http"

)

var (
	router *gin.Engine
)

func InitWebRouters() {
	router = gin.Default()

	// 用户组
	userRouter := router.Group("/project")
	{
		userRouter.GET("/deploy", project.Deploy)
		userRouter.GET("/display", project.Display)
	}
}

func InitHttpServer() {
	fmt.Println("Http Server 启动成功", initialize.ServerIp, initialize.HttpPort)
	http.ListenAndServe(":"+initialize.HttpPort, router)

}
