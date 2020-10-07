package node

import (
	"gaad/common"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
)

// @Description 创建集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/createCluster [post]
// @Tags 集群(Cluster)
func CreateCluster(c *gin.Context) {
	cluster := models.Cluster{}
	base.Create(c, &cluster, func(c *gin.Context) error {
		if cluster.ClusterName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// @Description 更新集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/updateCluster [put]
// @Tags 集群(Cluster)
func UpdateCluster(c *gin.Context) {

	cluster := models.Cluster{}
	base.Update(c, &cluster, func(c *gin.Context) error {
		if cluster.ClusterName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// @Description 为集群添加主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/setNode [post]
// @Tags 集群(Cluster)
func SetNode(c *gin.Context) {

	node := models.Node{}
	base.Update(c, &node, func(c *gin.Context) error {
		if node.ClusterId == 0 || node.NodeType == 0 {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		cluster := models.Cluster{}
		sqlitedb.First(&cluster, " id = ?", node.ClusterId)
		//cluster不存在
		if cluster.ID == 0 {
			return controllers.Response(c, common.ParameterIllegal, "所属于集群不存在", nil)
		}
		sqlitedb.First(&node, " id = ?", node.ID)
		node.ClusterId = cluster.ID
		node.ClusterName = cluster.ClusterName
		return nil
	})

}

// @Description 移除集群的主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/removeNode [delete]
// @Tags 集群(Cluster)
func RemoveNode(c *gin.Context) {

	node := models.Node{}

	err := base.GetModel(&node, c)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}
	if node.ID == 0 {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}

	sqlitedb.Delete(&node)
	//初始化关联集群数据
	node.ClusterId = 0
	node.ClusterName = ""
	node.NodeType = 0
	sqlitedb.Create(&node)
}

// @Description 查寻集群下面的所有主机
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/listNodes [post]
// @Tags 集群(Cluster)
func ListNodes(c *gin.Context) {
	node := models.Node{}
	var nodes []models.Node

	base.Page(c, &node, &nodes,
		func(c *gin.Context) error {
			if node.ClusterId == 0 {
				return controllers.Response(c, common.ParameterIllegal, "clusterId不能为空", nil)
			}
			return nil
		},
		func() (query interface{}, args []interface{}) {

			args = make([]interface{}, 3)

			sql := "1 = 1"
			if node.ClusterId != 0 {
				sql += " and cluster_id = "
				args = append(args, node.ClusterId)
			}
			query = sql
			return
		})

}

// @Description 删除集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/deleteCluster [delete]
// @Tags 集群(Cluster)
func DeleteCluster(c *gin.Context) {

	base.Delete(c, &models.Cluster{})
}

// @Description 分页查询集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/pageClusters [post]
// @Tags 集群(Cluster)
func PageClusters(c *gin.Context) {
	cluster := models.Cluster{}
	var clusters []models.Cluster

	base.Page(c, &cluster, &clusters,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 3)
			query = "cluster_name like ?"
			where[0] = "%" + cluster.ClusterName + "%"
			return
		})
}

// @Description 查寻主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/pageNodesForCluster [post]
// @Tags 集群(Cluster)
func PageNodesForCluster(c *gin.Context) {
	node := models.Node{}
	var nodes []models.Node

	base.Page(c, &node, &nodes,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, args []interface{}) {
			args = make([]interface{}, 3)

			sql := "1 = 1"
			sql += " and cluster_id = 0 and node_type =0"
			if node.Ip != "" {
				sql += " and ip like ?"
				args = append(args, "%"+node.Ip+"%")
			}

			query = sql
			return
		})
}

// @Description 查询所有集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} common.JsonResult
// @Router /cluster/listClusters [post]
// @Tags 集群(Cluster)
func ListClusters(c *gin.Context) {
	var clusters []models.Cluster
	base.List(c, &clusters, func() (where []interface{}) {
		where = make([]interface{}, 3)
		where[0] = "status = 1"
		return
	})
}
