package node

import (
	"gaad/common"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
)

// 查看全部在线用户
func CreateCluster(c *gin.Context) {
	cluster := models.Cluster{}
	base.Create(c, &cluster, func(c *gin.Context) error {
		if cluster.ClusterName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// 查看全部在线用户
func UpdateCluster(c *gin.Context) {

	cluster := models.Cluster{}
	base.Update(c, &cluster, func(c *gin.Context) error {
		if cluster.ClusterName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// 查看全部在线用户
func DeleteCluster(c *gin.Context) {

	base.Delete(c, &models.Cluster{})
}

// 查看全部在线用户
func PageClusters(c *gin.Context) {
	cluster := models.Cluster{}
	var clusters []models.Cluster

	base.Page(c, &cluster, &clusters, func() (query interface{}, where []interface{}) {
		where = make([]interface{}, 3)
		query = "cluster_name like ?"
		where[0] = "%" + cluster.ClusterName + "%"
		return
	})
}

func ListClusters(c *gin.Context) {
	var clusters []models.Cluster
	base.List(c, &clusters, func() (query interface{}, args []interface{}) {
		query = "status = 1"
		args = nil
		return
	})
}
