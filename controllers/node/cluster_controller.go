package node

import (
	"gaad/common"
	"gaad/common/web"
	"gaad/controllers"
	"gaad/db/sqlitedb"
	"gaad/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 查看全部在线用户
func CreateCluster(c *gin.Context) {

	modCluster("create", c)

}

// 查看全部在线用户
func UpdateCluster(c *gin.Context) {

	modCluster("update", c)

}

// 查看全部在线用户
func DeleteCluster(c *gin.Context) {

	cluster := models.Cluster{}
	web.GetModel(&cluster, c)

	sqlitedb.Delete(&cluster)
	controllers.Response(c, common.OK, "", nil)
}

func modCluster(operation string, c *gin.Context) {
	cluster := models.Cluster{}
	web.GetModel(&cluster, c)

	if cluster.ClusterName == "" {
		controllers.Response(c, common.ParameterIllegal, "", nil)
		return
	}
	if operation == "create" {
		sqlitedb.Create(&cluster)
	}
	if operation == "update" {
		sqlitedb.Update(&cluster)
	}

	controllers.Response(c, common.OK, "", nil)
}

// 查看全部在线用户
func PageClusters(c *gin.Context) {

	cluster := models.Cluster{}
	web.GetModel(&cluster, c)

	curPageStr := c.Query("curPage")
	pageRecordStr := c.Query("pageRecord")

	if curPageStr == "" {
		curPageStr = "1"
	}
	if pageRecordStr == "" {
		pageRecordStr = "10"
	}

	curPage, _ := strconv.Atoi(curPageStr)
	pageRecord, _ := strconv.Atoi(pageRecordStr)

	var (
		clusters []models.Cluster
	)

	total := sqlitedb.QueryPage(curPage, pageRecord, &clusters, "clusterName like ?", "%"+cluster.ClusterName+"%")

	data := make(map[string]interface{})

	data["clusters"] = clusters
	data["curPage"] = curPage
	data["pageRecord"] = pageRecord
	data["total"] = total
	controllers.Response(c, common.OK, "", data)
}

func ListClusters(c *gin.Context) {
	var (
		clusters []models.Cluster
	)
	sqlitedb.QueryList(clusters)

	data := make(map[string]interface{})
	data["clusters"] = clusters
	controllers.Response(c, common.OK, "", data)
}
