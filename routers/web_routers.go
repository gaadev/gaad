/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:20
 */

package routers

import (
	"fmt"
	"gaad/web/controllers/application"
	"gaad/web/controllers/node"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "gaad/docs"
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

	//处理跨域的问题
	corsMiddleWare := cors.Default()
	router.Use(corsMiddleWare)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		clusterRouter.POST("/setNode", node.SetNode)
		clusterRouter.POST("/listNodes", node.ListNodes)
		clusterRouter.DELETE("/removeNode", node.RemoveNode)
		clusterRouter.POST("/pageNodesForCluster", node.PageNodesForCluster)
	}

	projectRouter := router.Group("/project")
	{
		projectRouter.POST("/createProject", application.CreateProject)
		projectRouter.POST("/pageProjects", application.PageProjects)
		projectRouter.POST("/listProjects", application.ListProjects)
		projectRouter.PUT("/updateProject", application.UpdateProject)
		projectRouter.DELETE("/deleteProject", application.DeleteProject)
	}

	serviceRouter := router.Group("/service")
	{
		serviceRouter.POST("/deploy", application.Deploy)
		serviceRouter.POST("/display", application.Display)
	}

}

func InitHttpServer() {
	fmt.Println("Http Server 启动成功", initialize.ServerIp, initialize.HttpPort)
	http.ListenAndServe(":"+initialize.HttpPort, router)
}
