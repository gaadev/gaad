/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:20
 */

package routers

import (
	"fmt"
	"gaad/controllers/node"
	"gaad/controllers/project"
	"gaad/initialize"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	router *gin.Engine
)

func InitWebRouters() {
	router = gin.Default()

	corsMiddleWare := cors.Default()
	router.Use(corsMiddleWare)

	// 用户组
	projectRouter := router.Group("/project")
	{
		projectRouter.GET("/deploy", project.Deploy)
		projectRouter.GET("/display", project.Display)
	}
	nodeRouter := router.Group("/node")
	{
		nodeRouter.POST("/createNode", node.CreateNode)
		nodeRouter.POST("/pageNodes", node.PageNodes)
		nodeRouter.PUT("/updateNode", node.UpdateNode)
		nodeRouter.DELETE("/deleteNode", node.DeleteNode)
	}
	clusterRouter := router.Group("/cluster")
	{
		clusterRouter.POST("/createCluster", node.CreateCluster)
		clusterRouter.POST("/pageClusters", node.PageClusters)
		clusterRouter.POST("/listClusters", node.ListClusters)
		clusterRouter.PUT("/updateCluster", node.UpdateCluster)
		clusterRouter.DELETE("/deleteCluster", node.DeleteCluster)
	}
}

func InitHttpServer() {
	fmt.Println("Http Server 启动成功", initialize.ServerIp, initialize.HttpPort)
	http.ListenAndServe(":"+initialize.HttpPort, router)
}
