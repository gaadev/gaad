package node

import (
	"fmt"
	"gaad/common"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"net"
)

// @Description 创建主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/createNode [post]
func CreateNode(c *gin.Context) {
	node := models.Node{}
	base.Create(c, &node, func(c *gin.Context) error {
		address := net.ParseIP(node.Ip)
		if address == nil {
			return controllers.Response(c, common.ParameterIllegal, "Ip地址格式有误", nil)
		} else {
			fmt.Println("正确的ip地址", address.String())
		}

		if node.Port < 2 {
			return controllers.Response(c, common.ParameterIllegal, "端口不可以小于2", nil)
		}

		if node.Username == "" || node.Password == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// @Description 更新主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/updateNode [put]
func UpdateNode(c *gin.Context) {
	node := models.Node{}
	base.Update(c, &node, func(c *gin.Context) error {
		if node.Ip != "" {
			address := net.ParseIP(node.Ip)
			if address == nil {
				return controllers.Response(c, common.ParameterIllegal, "Ip地址格式有误", nil)
			} else {
				fmt.Println("正确的ip地址", address.String())
			}
		}

		if node.Port != 0 && node.Port < 2 {
			return controllers.Response(c, common.ParameterIllegal, "端口不可以小于2", nil)
		}
		return nil
	})

}

// @Description 删除主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/deleteNode [delete]
func DeleteNode(c *gin.Context) {

	base.Delete(c, &models.Node{})
}

// @Description 删除主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Param action query string false "excludeInCluster 排除已经加入集群的节点"
// @Success 200 {object} common.JsonResult
// @Router /node/pageNodes [post]
func PageNodes(c *gin.Context) {
	//查寻接口的行为
	action := c.Query("action")

	node := models.Node{}
	var nodes []models.Node

	base.Page(c, &node, &nodes,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, args []interface{}) {
			args = make([]interface{}, 3)

			sql := "1 = 1"
			if node.Ip != "" {
				sql += " and ip like ?"
				args = append(args, "%"+node.Ip+"%")
			}

			if action == "excludeInCluster" {
				sql += " and cluster_id > 0"
			}
			query = sql
			return
		})
}
