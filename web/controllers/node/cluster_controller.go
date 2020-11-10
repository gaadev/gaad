package node

import (
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/remote"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
)

// @Description 创建集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/createCluster [post]
// @Tags 集群(Cluster)
func CreateCluster(c *gin.Context) {
	cluster := models.Cluster{}
	rsp := base.Create(c, &cluster, func() *models.Rsp {
		if cluster.ClusterName == "" {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		return nil
	})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}

}

// @Description 更新集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/updateCluster [put]
// @Tags 集群(Cluster)
func UpdateCluster(c *gin.Context) {

	cluster := models.Cluster{}
	rsp := base.Update(c, &cluster, func() *models.Rsp {
		if cluster.ClusterName == "" {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		return nil
	})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 为集群添加主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/setNode [put]
// @Tags 集群(Cluster)
func SetNode(c *gin.Context) {

	node := models.Node{}
	rsp := base.Update(c, &node, func() *models.Rsp {
		if node.ClusterId == 0 || node.NodeType == 0 {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		cluster := models.Cluster{}
		sqlitedb.First(&cluster, " id = ?", node.ClusterId)
		//cluster不存在
		if cluster.ID == 0 {
			return controllers.Response(models.ParameterIllegal, "所属于集群不存在", nil)
		}
		//重新查询，防止该接口修改本接口功能之外的字段
		nodeOld := models.Node{}
		sqlitedb.First(&nodeOld, " id = ?", node.ID)
		nodeOld.ClusterId = cluster.ID
		nodeOld.ClusterName = cluster.ClusterName
		nodeOld.NodeType = node.NodeType
		node = nodeOld
		return nil
	})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 移除集群的主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/removeNode [delete]
// @Tags 集群(Cluster)
func RemoveNode(c *gin.Context) {

	node := models.Node{}

	err := base.GetModel(&node, c)
	if err != nil {
		controllers.Response(models.ParameterIllegal, "参数格式有误", nil).Write(c)
		return
	}
	if node.ID == 0 {
		controllers.Response(models.ParameterIllegal, "参数格式有误", nil).Write(c)
		return
	}

	sqlitedb.DeleteForce(&node)
	//初始化关联集群数据
	node.ClusterId = 0
	node.ClusterName = ""
	node.NodeType = 0

	sqlitedb.Create(&node)
	controllers.Response(models.OK, "", nil).Write(c)
}

// @Description 查寻集群下面的所有主机
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/listNodes [post]
// @Tags 集群(Cluster)
func ListNodes(c *gin.Context) {
	node := models.Node{}
	var nodes []models.Node

	rsp := base.Page(c, &node, &nodes,
		func() *models.Rsp {
			if node.ClusterId == 0 {
				return controllers.Response(models.ParameterIllegal, "clusterId不能为空", nil)
			}
			return nil
		},
		func() (query interface{}, args []interface{}) {

			args = make([]interface{}, 0)

			sql := "1 = 1"
			if node.ClusterId != 0 {
				sql += " and cluster_id = ?"
				args = append(args, node.ClusterId)
			}
			query = sql
			return
		})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 删除集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/deleteCluster [delete]
// @Tags 集群(Cluster)
func DeleteCluster(c *gin.Context) {

	cluster := models.Cluster{}
	err := base.GetModel(&cluster, c)
	if err != nil {
		controllers.Response(models.ParameterIllegal, "参数格式有误", nil).Write(c)
		return
	}
	var nodes []models.Node
	sqlitedb.QueryList(&nodes, "cluster_id = ?", cluster.ID)
	if len(nodes) > 0 {
		controllers.Response(models.OperationFailure, "请先删除集群所有子节点，再删除集群", nil).Write(c)
		return
	}

	sqlitedb.Delete(cluster)
	controllers.Response(models.OK, "", nil).Write(c)
}

// @Description 分页查询集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/pageClusters [post]
// @Tags 集群(Cluster)
func PageClusters(c *gin.Context) {
	cluster := models.Cluster{}
	var clusters []models.Cluster

	rsp := base.Page(c, &cluster, &clusters,
		func() *models.Rsp {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 0)

			sql := "1 = 1"
			if cluster.ClusterName != "" {
				query = "cluster_name like ?"
				where[0] = "%" + cluster.ClusterName + "%"
			}
			query = sql
			return
		})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 查寻主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/pageNodesForCluster [post]
// @Tags 集群(Cluster)
func PageNodesForCluster(c *gin.Context) {
	node := models.Node{}
	var nodes []models.Node

	rsp := base.Page(c, &node, &nodes,
		func() *models.Rsp {
			return nil
		},
		func() (query interface{}, args []interface{}) {
			args = make([]interface{}, 0)

			sql := "1 = 1"
			sql += " and cluster_id = 0 and node_type =0"
			if node.Ip != "" {
				sql += " and ip like ?"
				args = append(args, "%"+node.Ip+"%")
			}

			query = sql
			return
		})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 查询所有集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/listClusters [post]
// @Tags 集群(Cluster)
func ListClusters(c *gin.Context) {
	var clusters []models.Cluster
	rsp := base.List( &clusters, func() (where []interface{}) {
		where = make([]interface{}, 0)
		where = append(where, "status = 1")
		return
	})
	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 查询所有集群
// @Accept  json
// @Produce json
// @Param data body models.Cluster true "Data"
// @Success 200 {object} models.Rsp
// @Router /cluster/initCluster [post]
// @Tags 集群(Cluster)
func InitCluster(c *gin.Context) {
	cluster := models.Cluster{}

	err := base.GetModel(&cluster, c)
	if err != nil {
		controllers.Response(models.ParameterIllegal, "参数格式有误", nil).Write(c)
		return
	}
	clusterDb := models.Cluster{}
	sqlitedb.First(&clusterDb, "id = ? and status = 1", cluster.ID)
	if clusterDb.ID != 0 {
		var (
			nodes       []models.Node
			nodeMasters = make([]models.Node, 0)
			nodeSlavers = make([]models.Node, 0)
		)
		sqlitedb.QueryList(&nodes, "cluster_id = ?", clusterDb.ID)
		for _, node := range nodes {
			if node.NodeType == 2 {
				nodeMasters = append(nodeMasters, node)
			}
			if node.NodeType == 3 {
				nodeSlavers = append(nodeSlavers, node)
			}
		}

		if len(nodeMasters) < 1 || len(nodeSlavers) < 1 {
			controllers.Response(models.OperationFailure, "初始化集群时，至少存在一个master节点，一个slaver节点", nil).Write(c)
			return
		}

		switch clusterDb.Category {
		case "DockerSwarm":
			for i := 0; i < len(nodeMasters); i++ {
				if i == 0 {
					remote.InitDockerSwarmMaster(nodeMasters[i])
				} else {
					remote.FollowDockerSwarmMaster(nodeMasters[i])
				}

			}
			for i := 0; i < len(nodeSlavers); i++ {
				remote.InitDockerSwarmSlaver(nodeSlavers[i])
			}
		case "Kubernetes":
			for i := 0; i < len(nodeMasters); i++ {
				if i == 0 {
					remote.InitKubernetesMaster(nodeMasters[i])
				} else {
					remote.FollowKubernetesMaster(nodeMasters[i])
				}

			}
			for i := 0; i < len(nodeSlavers); i++ {
				remote.InitKubernetesSlaver(nodeSlavers[i])
			}
		default:
			controllers.Response(models.OperationFailure, "集群类型不支持", nil).Write(c)
			return
		}

	}

	controllers.Response(models.OK, "", nil).Write(c)
}
